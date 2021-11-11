package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/friendsofgo/errors"
	"github.com/rogpeppe/go-internal/semver"
	"github.com/spf13/cobra"

	"github.com/gmhafiz/go8gen/internal/app"
)

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&DirPath, "path", "p", "", "path to new project")
}

var (
	DirPath = ""
)

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

		a := app.New()

		name := strings.ToLower(args[0])
		if len(args) == 2 {
			a.Project.ModuleName = args[1]
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			log.Fatal(err)
		}

		path, name = parseProjectNameAndPath(path, name)
		a.Project.Path = path
		a.Project.Name = name

		a.Project.Address = "http://localhost:3080/api/v1"

		if a.IsDirectoryExists(a.Project.Path) {
			log.Fatalf("chosen directory name already exists: %s", a.Project.Name)
		}

		a.Prompt()

		directories := []string{
			"cmd",
			"cmd/migrate",
			"cmd/migrate/migrate",
			"configs",
			"database",
			"database/migrations",
			"internal",
			"internal/domain",
			"internal/domain/health",
			"internal/domain/health/handler",
			"internal/domain/health/handler/http",
			"internal/domain/health/repository/database",
			"internal/domain/health/usecase",
			"internal/middleware",
			"internal/models",
			"internal/server",
			"internal/utility",
			"internal/utility/database",
			"internal/utility/filter",
			"internal/utility/message",
			"internal/utility/param",
			"internal/utility/respond",
			"internal/utility/validate",
			"scripts",
			"third_party",
			"third_party/database",
			"third_party/validate",
		}
		structure := []app.Structure{
			{
				TemplateFileName: ".air.toml.tmpl",
				FileName:         ".air.toml",
				Parse:            true,
			},
			{
				TemplateFileName: ".env.tmpl",
				FileName:         ".env",
				Parse:            true,
			},
			{
				TemplateFileName: ".gitignore.tmpl",
				FileName:         ".gitignore",
				Parse:            false,
			},
			{
				TemplateFileName: "Taskfile.tmpl",
				FileName:         "Taskfile.yml",
				Parse:            false,
			},
			{
				TemplateFileName: "cmd/migrate/main.go.tmpl",
				FileName:         "cmd/migrate/main.go",
				Parse:            true,
			},
			{
				TemplateFileName: "cmd/migrate/migrate/migrate.go.tmpl",
				FileName:         "cmd/migrate/migrate/migrate.go",
				Parse:            true,
			},
			{
				TemplateFileName: "cmd/go8/main.go.tmpl",
				FileName:         fmt.Sprintf("cmd/%s/main.go", a.Project.Name),
				Parse:            true,
			},
			{
				TemplateFileName: "configs/api.go.tmpl",
				FileName:         "configs/api.go",
				Parse:            true,
			},
			{
				TemplateFileName: "configs/cache.go.tmpl",
				FileName:         "configs/cache.go",
				Parse:            true,
			},
			{
				TemplateFileName: "configs/configs.go.tmpl",
				FileName:         "configs/configs.go",
				Parse:            true,
			},
			{
				TemplateFileName: "configs/database.go.tmpl",
				FileName:         "configs/database.go",
				Parse:            true,
			},
			{
				TemplateFileName: "configs/elasticsearch.go.tmpl",
				FileName:         "configs/elasticsearch.go",
				Parse:            true,
			},
			{
				TemplateFileName: "server/server.go.tmpl",
				FileName:         "internal/server/server.go",
				Parse:            true,
			},
			{
				TemplateFileName: "server/initDomains.go.tmpl",
				FileName:         "internal/server/initDomains.go",
				Parse:            true,
			},
			{
				TemplateFileName: "health/handler.go.tmpl",
				FileName:         "internal/domain/health/handler.go",
				Parse:            false,
			},
			{
				TemplateFileName: "health/http/handler.go.tmpl",
				FileName:         "internal/domain/health/handler/http/handler.go",
				Parse:            true,
			},
			{
				TemplateFileName: "health/http/register.go.tmpl",
				FileName:         "internal/domain/health/handler/http/register.go",
				Parse:            true,
			},
			{
				TemplateFileName: "health/database/database.go.tmpl",
				FileName:         "internal/domain/health/repository/database/database.go",
				Parse:            true,
			},
			{
				TemplateFileName: "health/usecase/usecase.go.tmpl",
				FileName:         "internal/domain/health/usecase/usecase.go",
				Parse:            true,
			},
			{
				TemplateFileName: "health/usecase.go.tmpl",
				FileName:         "internal/domain/health/usecase.go",
				Parse:            true,
			},
			{
				TemplateFileName: "health/repository.go.tmpl",
				FileName:         "internal/domain/health/repository.go",
				Parse:            true,
			},
			{
				TemplateFileName: "middleware/cors.go.tmpl",
				FileName:         "internal/middleware/cors.go",
				Parse:            true,
			},
			{
				TemplateFileName: "middleware/json.go.tmpl",
				FileName:         "internal/middleware/json.go",
				Parse:            true,
			},
			{
				TemplateFileName: "third_party/database/sqlx.go.tmpl",
				FileName:         "third_party/database/sqlx.go",
				Parse:            true,
			},
			{
				TemplateFileName: "third_party/validate/validate.go.tmpl",
				FileName:         "third_party/validate/validate.go",
				Parse:            false,
			},
			{
				TemplateFileName: "utility/database/db.go.tmpl",
				FileName:         "internal/utility/database/db.go",
				Parse:            true,
			},
			{
				TemplateFileName: "utility/filter/base.go.tmpl",
				FileName:         "internal/utility/filter/base.go",
				Parse:            true,
			},
			{
				TemplateFileName: "utility/respond/error.go.tmpl",
				FileName:         "internal/utility/respond/error.go",
				Parse:            true,
			},
			{
				TemplateFileName: "utility/respond/render.go.tmpl",
				FileName:         "internal/utility/respond/render.go",
				Parse:            true,
			},
			{
				TemplateFileName: "utility/param/param.go.tmpl",
				FileName:         "internal/utility/param/param.go",
				Parse:            true,
			},
			{
				TemplateFileName: "utility/validate/validate.go.tmpl",
				FileName:         "internal/utility/validate/validate.go",
				Parse:            true,
			},
			{
				TemplateFileName: "scripts/install-task.sh",
				FileName:         "scripts/install-task.sh",
				Parse:            false,
			},
			{
				TemplateFileName: "scripts/stopDockertestByPort.sh",
				FileName:         "scripts/stopDockertestByPort.sh",
				Parse:            false,
			},
			{
				TemplateFileName: "README.md.tmpl",
				FileName:         "README.md",
				Parse:            false,
			},
			{
				TemplateFileName: "docker-compose.yml.tmpl",
				FileName:         "docker-compose.yml",
				Parse:            false,
			},
			{
				TemplateFileName: "Dockerfile.tmpl",
				FileName:         "Dockerfile",
				Parse:            false,
			},
			{
				TemplateFileName: "sqlboiler.toml.tmpl",
				FileName:         "sqlboiler.toml",
				Parse:            true,
			},
		}

		a.SetStructure(structure)

		syscall.Umask(0)
		err = os.Mkdir(a.Project.Path, 0755)
		if err != nil {
			a.Fatal(errors.Wrap(err, "error creating directory: %s"), a.Project.Name)

		}
		err = os.Chdir(a.Project.Path)
		if err != nil {
			a.Fatal(errors.Wrap(err, "error going into directory: %s"), a.Project.Name)
		}

		projectNamePath := fmt.Sprintf("cmd/%s", a.Project.Name)

		directories = append(directories, projectNamePath)
		err = a.CreateDirectories(directories)
		if err != nil {
			a.Fatal(err, a.Project.Name)
		}

		err = a.CreateFiles()
		if err != nil {
			a.Fatal(err)
		}

		err = a.InitGoMod()
		if err != nil {
			a.Fatal(err)
		}

		fmt.Printf(InfoColor, "...done.\n")
		fmt.Println("\nChange directory to")
		fmt.Printf(InfoColor, fmt.Sprintf("    cd %s\n", a.Project.Path))
		fmt.Println("You may now run the following command to download dependencies")
		fmt.Printf(InfoColor, "    go mod tidy\n")
		fmt.Println("Run the API with")
		fmt.Printf(InfoColor, fmt.Sprintf("    go run cmd/%s/main.go\n\n", a.Project.ModuleName))
	},
}

func parseProjectNameAndPath(path, name string) (projectPath, projectName string) {
	if path != "" {
		_, projectName = filepath.Split(path)
		projectPath = path
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		projectPath = filepath.Join(cwd, name)
		projectName = name
	}

	return projectPath, projectName
}

func checkVersion() {
	version := runtime.Version()
	res := semver.Compare(version, "1.14")
	if res < 0 {
		log.Fatal("go version must be > 1.14 for go modules support")
	}

	goSupport := semver.Compare(version, "1.16")
	if goSupport < 0 {
		log.Println("warning: only last 2 versions of go are supported officially by the Go team. See https://golang.org/doc/devel/release#policy")
	}
}
