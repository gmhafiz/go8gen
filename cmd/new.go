package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	directories []string
	files       []Structure
)

type Structure struct {
	TemplateFileName string
	FileName         string
	Parse            bool
}

type Semver struct {
	Program string
	Major   int
	Minor   int
	Patch   int
}

func init() {
	rootCmd.AddCommand(newCmd)

	directories = []string{
		"cmd",
		"configs",
		"database",
		"internal",
		"internal/domain",
		"internal/domain/health",
		"internal/domain/health/handler",
		"internal/domain/health/handler/http",
		"internal/domain/health/repository/postgres",
		"internal/domain/health/usecase",
		"internal/middleware",
		"internal/model",
		"internal/resource",
		"internal/server",
		"scripts",
		"third_party",
		"third_party/database",
	}
	files = []Structure{}
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new starter",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		checkVersion()

		if len(args) == 0 {
			log.Fatal("must put a project name. e.g. : go8 new new_project")
		}
		projectName := args[0]

		projectNamePath := fmt.Sprintf("cmd/%s", projectName)
		directories = append(directories, projectNamePath)
		err := createDirectories(directories, projectName)
		if err != nil {
			log.Fatal(err)
		}
		err = initGoMod(projectName)
		if err != nil {
			log.Fatal(err)
		}

		structures := []Structure{
			{
				TemplateFileName: "../tmpl/.env.tmpl",
				FileName:         fmt.Sprint(".env"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/.gitignore.tmpl",
				FileName:         fmt.Sprint(".gitignore"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/Taskfile.tmpl",
				FileName:         fmt.Sprint("Taskfile.yml"),
				Parse:            false,
			},
			{
				TemplateFileName: "../tmpl/cmd/main.tmpl",
				FileName:         fmt.Sprintf("cmd/%s/%s.go", projectName, projectName),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/configs/api.go.tmpl",
				FileName:         fmt.Sprint("configs/api.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/configs/cache.go.tmpl",
				FileName:         fmt.Sprint("configs/cache.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/configs/configs.go.tmpl",
				FileName:         fmt.Sprint("configs/configs.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/configs/database.go.tmpl",
				FileName:         fmt.Sprint("configs/database.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/configs/elasticsearch.go.tmpl",
				FileName:         fmt.Sprint("configs/elasticsearch.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/server/server.go.tmpl",
				FileName:         fmt.Sprint("internal/server/server.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/health/http/handler.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/handler/http/handler.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/health/http/register.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/handler/http/register.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/health/postgres/postgres.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/repository/postgres/postgres.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/health/usecase/usecase.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/usecase/usecase.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/health/usecase.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/usecase.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/health/repository.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/repository.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/middleware/cors.go.tmpl",
				FileName:         fmt.Sprintf("internal/middleware/cors.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/middleware/paginate.go.tmpl",
				FileName:         fmt.Sprintf("internal/middleware/paginate.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/third_party/database/database.go.tmpl",
				FileName:         fmt.Sprintf("third_party/database/database.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "../tmpl/scripts/install-task.sh",
				FileName:         fmt.Sprintf("scripts/install-task.sh"),
				Parse:            false,
			},
			{
				TemplateFileName: "../tmpl/README.md.tmpl",
				FileName:         fmt.Sprintf("README.md"),
				Parse:            false,
			},
		}
		err = createFiles(projectName, "", structures)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(InfoColor, "...done.\n")
	},
}

func checkVersion() {
	version := runtime.Version()
	parsed := parseVersion(version)
	if parsed.Major < 1 && parsed.Minor < 13 {
		log.Fatal(errors.New("go version must be > 1.13 for go modules support"))
	}
}

func parseVersion(version string) *Semver {
	splitted := strings.Split(version, ".")
	major, err := strconv.Atoi(splitted[0][2:])
	if err != nil {
		log.Fatal("error parsing Go version")
	}
	minor, err := strconv.Atoi(splitted[1])
	if err != nil {
		log.Fatal("error parsing Go version")
	}
	patch, err := strconv.Atoi(splitted[2])
	if err != nil {
		log.Fatal("error parsing Go version")
	}
	return &Semver{
		Program: "go",
		Major:   major,
		Minor:   minor,
		Patch:   patch,
	}
}

func createDirectories(directories []string, projectName string) error {
	for _, val := range directories {
		syscall.Umask(0)
		err := os.MkdirAll(val, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func initGoMod(projectName string) error {
	cmd := exec.Command("go", "mod", "init", projectName)
	_, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "error running: go mod init")
	}
	return nil
}

func createFiles(projectName, domain string, structures []Structure) error {
	vars := struct {
		ProjectName     string
		Domain          string
		DomainLowerCase string
	}{
		ProjectName:     projectName,
		Domain:          strings.Title(domain),
		DomainLowerCase: domain,
	}

	for _, val := range structures {
		files = append(files, val)
	}

	for _, val := range files {
		err := createFile(val.TemplateFileName, val.FileName, vars, val.Parse)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFile(tmplFileName, fileName string, vars struct{ ProjectName, Domain, DomainLowerCase string },
parse bool) error {
	file, err := os.Create(fileName)
	if err != nil {
		return errors.Wrapf(err, "error creating file: %s", fileName)
	}

	if !parse {
		return copyFile(fileName, tmplFileName)
	} else {
		return parseFile(file, tmplFileName, vars)
	}
}

func parseFile(file *os.File, tmplFileName string, vars struct{ ProjectName, Domain, DomainLowerCase string }) error {
	tmpl, err := template.ParseFiles(tmplFileName)
	if err != nil {
		return errors.Wrapf(err, "error parsing file: %s", tmplFileName)
	}
	return tmpl.Execute(file, vars)
}

func copyFile(fileName, tmplFileName string) error {
	tmplContent, err := ioutil.ReadFile(tmplFileName)
	if err != nil {
		return errors.Wrapf(err, "error reading file: %s", tmplFileName)
	}
	fCreate, err := os.Create(fileName)
	if err != nil {
		return errors.Wrapf(err, "error opening file: %s", tmplFileName)
	}
	_, err = fCreate.WriteString(string(tmplContent))
	if err != nil {
		return errors.Wrapf(err, "error writing content into: %s", fileName)
	}
	return nil
}
