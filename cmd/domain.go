package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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
			log.Fatal("must put a name for the domain. e.g. : go8 domain book")
		}

		a := app.New()
		p := &app.Project{
			Name:                   getProjectName(),
			ModuleName:             getModuleName(),
			Domain:                 strings.Title(args[0]),
			DomainLowerCase:        strings.ToLower(args[0]),
			ScaffoldAuthentication: false,
			ScaffoldUseCase:        true,
			ScaffoldRepository:     true,
		}
		if p.ModuleName == "" {
			log.Fatal("error finding module name")
		}
		a.SetProject(p)

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
				TemplateFileName: "../tmpl/domain/repository.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/repository.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/repository/postgres/postgres.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/repository/postgres/postgres.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/repository/postgres/postgres_test.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/repository/postgres/postgres_test.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/usecase.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/usecase.go", p.DomainLowerCase),
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
				TemplateFileName: "../tmpl/domain/handler.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/handler.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/http/handler.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/handler/http/handler.go", p.DomainLowerCase),
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
		fmt.Sprintf("internal/domain/%s/repository/postgres", domain),
		fmt.Sprintf("internal/domain/%s/usecase", domain),
	}

	return directories
}
