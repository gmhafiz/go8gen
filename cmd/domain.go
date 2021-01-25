package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

func init() {
	rootCmd.AddCommand(domainCmd)

	directories = []string{}
	files = []Structure{}
	p = Project{}
}

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Create a new domain including controller, use case, and repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("must put a name for the domain. e.g. : go8 domain book")
		}

		p := Project{}
		p.Domain = strings.Title(args[0])
		p.DomainLowerCase = strings.ToLower(args[0])
		p.Name = getProjectName()
		p.ModuleName = getModuleName()
		if p.ModuleName == "" {
			log.Fatal("error finding module name")
		}

		directories := createDirectoryNames(p.DomainLowerCase)
		err := createDirectories(directories)
		if err != nil {
			log.Fatal(err)
		}

		structures := []Structure{
			{
				TemplateFileName: "../tmpl/domain/repository.go.tmpl",
				FileName: fmt.Sprintf("internal/domain/%s/repository.go", p.DomainLowerCase),
				Parse: true,
			},
			{
				TemplateFileName: "../tmpl/domain/repository/postgres/postgres.go.tmpl",
				FileName: fmt.Sprintf("internal/domain/%s/repository/postgres/postgres.go", p.DomainLowerCase),
				Parse: true,
			},
			{
				TemplateFileName: "../tmpl/domain/usecase.go.tmpl",
				FileName: fmt.Sprintf("internal/domain/%s/usecase.go", p.DomainLowerCase),
				Parse: true,
			},
			{
				TemplateFileName: "../tmpl/domain/usecase/usecase.go.tmpl",
				FileName: fmt.Sprintf("internal/domain/%s/usecase/usecase.go", p.DomainLowerCase),
				Parse: true,
			},
			{
				TemplateFileName: "../tmpl/domain/http/handler.go.tmpl",
				FileName: fmt.Sprintf("internal/domain/%s/handler/http/handler.go", p.DomainLowerCase),
				Parse: true,
			},
			{
				TemplateFileName: "../tmpl/domain/http/register.go.tmpl",
				FileName: fmt.Sprintf("internal/domain/%s/handler/http/register.go", p.DomainLowerCase),
				Parse: true,
			},
			{
				TemplateFileName: "../tmpl/model/model.go.tmpl",
				FileName: fmt.Sprintf("internal/model/%s.go", p.DomainLowerCase),
				Parse: true,
			},
			{
				TemplateFileName: "../tmpl/resource/resource.go.tmpl",
				FileName: fmt.Sprintf("internal/resource/%s.go", p.DomainLowerCase),
				Parse: true,
			},
		}
		err = createFiles(p, structures)
		if err != nil {
			log.Fatalf(ErrorColor, err)
		}

		err = injectCode(p)
		if err != nil {
			log.Fatalf(ErrorColor, err)
		}

		fmt.Printf(InfoColor, "...done.\n")
	},
}

func injectCode(p Project) error {
	const serverFileName = "internal/server/server.go"
	const injectImport = "// inject:import"
	const injectApp = "//inject:app"
	const injectUseCase = "// inject:usecase"
	const injectHandler = "// inject:handler"
	p.Domain = strings.Title(p.DomainLowerCase)
	importTmpl1 := fmt.Sprintf(`%sHTTP "%s/internal/domain/%s/handler/http"`, p.DomainLowerCase, p.Name, p.DomainLowerCase)
	importTmpl2 := fmt.Sprintf(`%sPostgres "%s/internal/domain/%s/repository/postgres"`, p.DomainLowerCase, p.Name, p.DomainLowerCase)
	importTmpl3 := fmt.Sprintf(`%sUseCase "%s/internal/domain/%s/usecase"`, p.DomainLowerCase, p.Name, p.DomainLowerCase)
	appTmpl := fmt.Sprintf(`%sUC *%sUseCase.%sUseCase`, p.DomainLowerCase, p.DomainLowerCase, p.Domain)
	usecaseTmpl := fmt.Sprintf(`%sUC: %sUseCase.New%sUseCase(%sPostgres.New%sRepository(db)),`, p.DomainLowerCase, p.DomainLowerCase, p.Domain, p.DomainLowerCase, p.Domain)
	handlerTmpl := fmt.Sprintf(`%sHTTP.RegisterHTTPEndPoints(router, a.%sUC)`, p.DomainLowerCase, p.DomainLowerCase)

	serverContent, err := ioutil.ReadFile(serverFileName)
	if err != nil {
		return errors.Wrapf(err, "error reading file: ", serverFileName)
	}

	var newFile []string
	temp := strings.Split(string(serverContent), "\n")
	for _, line := range temp {
		newFile = append(newFile, line)
		stripped := strings.Trim(line, "\t")
		stripped = strings.Trim(stripped, "\n")
		if stripped == injectImport {
			newFile = append(newFile, importTmpl1)
			newFile = append(newFile, importTmpl2)
			newFile = append(newFile, importTmpl3)
		}
		if stripped == injectApp {
			newFile = append(newFile, appTmpl)
		}
		if stripped == injectUseCase {
			newFile = append(newFile, usecaseTmpl)
		}
		if stripped == injectHandler {
			newFile = append(newFile, handlerTmpl)
		}
	}

	fCreate, err := os.Create(serverFileName)
	if err != nil {
		return errors.Wrapf(err, "error opening file: %s", serverFileName)
	}
	_, err = fCreate.WriteString(strings.Join(newFile, "\n"))
	if err != nil {
		return errors.Wrapf(err, "error writing file: %s", serverFileName)
	}
	return nil
}

// adapted from https://stackoverflow.com/a/63393712/1033134
func getProjectName() string {
	goModBytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		os.Exit(0)
	}
	return modfile.ModulePath(goModBytes)
}

func getModuleName() string {
	currPath, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting current path")
	}
	fileWithGoModule := "internal/domain/health/handler/http/handler.go"
	file, err := ioutil.ReadFile(filepath.Join(currPath, fileWithGoModule))
	if err != nil {
		log.Fatalf("error finding the file with go modules name")
	}
	temp := strings.Split(string(file), "\n")
	for _, line := range temp {
		stripped := strings.Trim(line, "\t")
		stripped = strings.Trim(stripped, "\n")
		if strings.Contains(stripped, "internal/domain/health") {
			return strings.Trim(strings.Split(stripped, "/")[0], "\"")
		}
	}
	return ""
}

func createDirectoryNames(domain string) []string {
	directories := []string{
		fmt.Sprintf("internal/domain/%s", domain),
		fmt.Sprintf("internal/domain/%s/handler", domain),
		fmt.Sprintf("internal/domain/%s/handler/http", domain),
		fmt.Sprintf("internal/domain/%s/repository", domain),
		fmt.Sprintf("internal/domain/%s/repository/postgres", domain),
		fmt.Sprintf("internal/domain/%s/usecase", domain),
	}

	return directories
}
