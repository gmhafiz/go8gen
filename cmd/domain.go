package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"

	"github.com/gmhafiz/go8gen/internal/app"
)

func init() {
	rootCmd.AddCommand(domainCmd)
}

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Create a new domain including controller, use case, and repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Must put a name for the domain. e.g. : go8 domain book")
		}

		a := app.New()
		p := &app.Project{
			ModuleName:             getProjectName(),
			Domain:                 strings.Title(args[0]),
			DomainLowerCase:        strings.ToLower(args[0]),
			ScaffoldAuthentication: false,
			ScaffoldUseCase:        true,
			ScaffoldRepository:     true,
			ExpandedID:             true,
			Database:               app.Database{Type: getDatabaseType()},
		}
		if p.ModuleName == "" {
			log.Fatal("error finding module name")
		}
		a.SetProject(p)

		if domainExists(p.Domain) {
			log.Fatal("domain has already exists")
		}

		directories := createDirectoryNames(p.DomainLowerCase)
		err := a.CreateDirectories(directories)
		if err != nil {
			log.Fatal(err)
		}

		structure := []app.Structure{
			{
				TemplateFileName: "examples/domain.http.tmpl",
				FileName:         fmt.Sprintf("examples/%s.http", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/repository/database/database.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/repository/database/database.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/repository/database/database_test.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/repository/database/database_test.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/usecase/usecase.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/usecase/usecase.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/usecase/usecase_test.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/usecase/usecase_test.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/http/handler_test.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/handler/http/handler_test.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/http/register.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/handler/http/register.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/filters.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/filters.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/resource.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/resource.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/request.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/request.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/models/model.go.tmpl",
				FileName:         fmt.Sprintf("internal/models/%s.go", p.DomainLowerCase),
				Parse:            true,
			},
		}
		a.SetStructure(structure)

		err = a.CreateFiles()
		if err != nil {
			log.Fatalf(ErrorColor, err)
		}

		err = p.InjectImportDomainHandlerCode()
		if err != nil {
			log.Fatalf(ErrorColor, err)
		}

		err = a.GoModTidy()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(InfoColor, "...done.\n")
	},
}

func domainExists(domain string) bool {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pathToDomain := filepath.Join(dir, "internal", "domain", strings.ToLower(domain))
	_, err = os.Stat(pathToDomain)
	if err != nil {
		return false
	}

	return true
}

func getDatabaseType() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	type dbType struct {
		Driver string `envconfig:"DB_DRIVER"`
	}
	var t dbType
	err = envconfig.Process("DB", &t)
	if err != nil {
		log.Fatal(err)
	}

	return t.Driver
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
		"examples",
		fmt.Sprintf("internal/domain/%s", domain),
		fmt.Sprintf("internal/domain/%s/handler", domain),
		fmt.Sprintf("internal/domain/%s/handler/http", domain),
		fmt.Sprintf("internal/domain/%s/repository", domain),
		fmt.Sprintf("internal/domain/%s/repository/database", domain),
		fmt.Sprintf("internal/domain/%s/usecase", domain),
	}

	return directories
}
