package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"text/template"

	"github.com/lu4p/binclude"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var tmplPath = binclude.Include("../tmpl")

var (
	directories []string
	files       []Structure
	p           Project
)

type Project struct {
	Name            string
	ModuleName      string
	Domain          string
	DomainLowerCase string
}

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

	p = Project{}
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new starter",
	Long:  ``,
	Example: `1. go8 new <projectName>
2. go8 new <projectName> <site.com/user/repoName> # optional`,
	Run: func(cmd *cobra.Command, args []string) {
		checkVersion()

		if len(args) == 0 {
			log.Fatal("must put a name. Example:\ngo8 new new_project")
		}

		p.Name = strings.ToLower(args[0])
		if len(args) == 1 {
			p.ModuleName = args[0]
		} else if len(args) == 2 {
			p.ModuleName = args[1]
		}

		if isDirectoryExists(p.Name) {
			log.Fatalf("chosen directory name already exists: %s", p.Name)
		}

		syscall.Umask(0)
		err := os.Mkdir(p.Name, 0755)
		if err != nil {
			Fatal(errors.Wrap(err, "error creating directory: %s"), p.Name)

		}
		err = os.Chdir(p.Name)
		if err != nil {
			Fatal(errors.Wrap(err, "error going into directory: %s"), p.Name)
		}

		projectNamePath := fmt.Sprintf("cmd/%s", p.Name)
		directories = append(directories, projectNamePath)
		err = createDirectories(directories)
		if err != nil {
			Fatal(err, p.Name)
		}
		err = initGoMod(p.Name)
		if err != nil {
			Fatal(err)
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
				FileName:         fmt.Sprintf("cmd/%s/%s.go", p.Name, p.Name),
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
			{
				TemplateFileName: "../tmpl/docker-compose.yml.tmpl",
				FileName:         fmt.Sprintf("docker-compose.yml"),
				Parse:            false,
			},
			{
				TemplateFileName: "../tmpl/Dockerfile.tmpl",
				FileName:         fmt.Sprintf("Dockerfile"),
				Parse:            false,
			},
		}
		err = createFiles(p, structures)
		if err != nil {
			Fatal(err)
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

func isDirectoryExists(projectName string) bool {
	_, err := os.Stat(projectName)
	if err != nil {
		return false
	}
	return true
}

func Fatal(msg error, args ...string) {
	err := os.RemoveAll(p.Name)
	if err != nil {
		log.Printf("error removing diectory: %s", p.Name)
	}
	log.Fatalf(msg.Error(), args)
}

func createDirectories(directories []string) error {
	for _, val := range directories {
		syscall.Umask(0)
		err := os.MkdirAll(val, 0755)
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

func createFiles(project Project, structures []Structure) error {
	for _, val := range structures {
		files = append(files, val)
	}

	for _, val := range files {
		err := createFile(val.TemplateFileName, val.FileName, project, val.Parse)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFile(tmplFileName, fileName string, project Project, parse bool) error {
	file, err := os.Create(fileName)
	if err != nil {
		return errors.Wrapf(err, "error creating file: %s", fileName)
	}

	if !parse {
		return copyFile(fileName, tmplFileName)
	} else {
		return parseFile(file, tmplFileName, project)
	}
}

func parseFile(file *os.File, tmplFileName string, project Project) error {
	tmplContent, err := BinFS.ReadFile(tmplFileName)
	fileNameTail := filepath.Base(tmplFileName)
	f, err := os.Create(filepath.Join("/tmp/", fileNameTail))
	if err != nil {
		return errors.Wrap(err, "error creating file at /tmp folder")
	}
	_, err = f.Write(tmplContent)
	if err != nil {
		return errors.Wrap(err, "error writing temporary file")
	}

	tmpl, err := template.ParseFiles(filepath.Join("/tmp", fileNameTail))

	if err != nil {
		return errors.Wrapf(err, "error parsing file: %s", fileNameTail)
	}
	return tmpl.Execute(file, project)
}

func copyFile(fileName, tmplFileName string) error {
	filePath := filepath.Join(tmplPath, tmplFileName)
	tmplContent, err := BinFS.ReadFile(filePath)
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
