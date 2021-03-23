package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gmhafiz/go8gen/internal/app"
)

func init() {
	rootCmd.AddCommand(handlerCmd)
}

var handlerCmd = &cobra.Command{
	Use:   "handler",
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

		directories := createHandlerDirectoryNames(p.DomainLowerCase, p.ScaffoldUseCase, p.ScaffoldRepository)
		err := a.CreateDirectories(directories)
		if err != nil {
			log.Fatal(err)
		}

		structure := []app.Structure{
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
				TemplateFileName: "../tmpl/domain/usecase.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/usecase.go", p.DomainLowerCase),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/domain/usecase/usecase.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/%s/usecase/usecase.go", p.DomainLowerCase),
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

		fmt.Printf(InfoColor, "...done.\n")

		err = a.GoModTidy()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func createHandlerDirectoryNames(domain string, useCase, repo bool) []string {
	directories := []string{
		fmt.Sprintf("internal/domain/%s", domain),
		fmt.Sprintf("internal/domain/%s/handler", domain),
		fmt.Sprintf("internal/domain/%s/handler/http", domain),
	}
	if useCase {
		directories = append(directories, fmt.Sprintf("internal/domain/%s/usecase/http", domain))
	}
	return directories
}
