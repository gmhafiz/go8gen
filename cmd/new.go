package cmd

import (
	"fmt"
	"log"
	"os"
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

		a := app.New()

		a.Project.Name = strings.ToLower(args[0])
		if len(args) == 2 {
			a.Project.ModuleName = args[1]
		}

		a.Prompt()

		directories := []string{
			"cmd",
			"configs",
			"database",
			"database/migrations",
			"internal",
			"internal/domain",
			"internal/domain/health",
			"internal/domain/health/handler",
			"internal/domain/health/handler/http",
			"internal/domain/health/repository/postgres",
			"internal/domain/health/usecase",
			"internal/middleware",
			"internal/model",
			"internal/server",
			"internal/utility",
			"internal/utility/database",
			"internal/utility/respond",
			"scripts",
			"third_party",
			"third_party/database",
		}
		structure := []app.Structure{
			{
				TemplateFileName: ".env.tmpl",
				FileName:         fmt.Sprint(".env"),
				Parse:            true,
			},
			{
				TemplateFileName: ".gitignore.tmpl",
				FileName:         fmt.Sprint(".gitignore"),
				Parse:            false,
			},
			{
				TemplateFileName: "Taskfile.tmpl",
				FileName:         fmt.Sprint("Taskfile.yml"),
				Parse:            false,
			},
			{
				TemplateFileName: "cmd/main.go.tmpl",
				FileName:         fmt.Sprintf("cmd/%s/main.go", a.Project.Name),
				Parse:            true,
			},
			{
				TemplateFileName: "configs/api.go.tmpl",
				FileName:         fmt.Sprint("configs/api.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "configs/cache.go.tmpl",
				FileName:         fmt.Sprint("configs/cache.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "configs/configs.go.tmpl",
				FileName:         fmt.Sprint("configs/configs.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "configs/database.go.tmpl",
				FileName:         fmt.Sprint("configs/database.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "configs/elasticsearch.go.tmpl",
				FileName:         fmt.Sprint("configs/elasticsearch.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "server/server.go.tmpl",
				FileName:         fmt.Sprint("internal/server/server.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "health/http/handler.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/handler/http/handler.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "health/http/register.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/handler/http/register.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "health/postgres/postgres.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/repository/postgres/postgres.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "health/usecase/usecase.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/usecase/usecase.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "health/usecase.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/usecase.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "health/repository.go.tmpl",
				FileName:         fmt.Sprintf("internal/domain/health/repository.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "middleware/cors.go.tmpl",
				FileName:         fmt.Sprintf("internal/middleware/cors.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "middleware/json.go.tmpl",
				FileName:         fmt.Sprintf("internal/middleware/json.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "middleware/paginate.go.tmpl",
				FileName:         fmt.Sprintf("internal/middleware/paginate.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "third_party/database/sqlx.go.tmpl",
				FileName:         fmt.Sprintf("third_party/database/sqlx.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "utility/database/db.go.tmpl",
				FileName:         fmt.Sprintf("internal/utility/database/db.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "utility/respond/render.go.tmpl",
				FileName:         fmt.Sprintf("internal/utility/respond/render.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "utility/respond/urlParam.go.tmpl",
				FileName:         fmt.Sprintf("internal/utility/respond/urlParam.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "utility/respond/validate.go.tmpl",
				FileName:         fmt.Sprintf("internal/utility/respond/validate.go"),
				Parse:            true,
			},
			{
				TemplateFileName: "scripts/install-task.sh",
				FileName:         fmt.Sprintf("scripts/install-task.sh"),
				Parse:            false,
			},
			{
				TemplateFileName: "README.md.tmpl",
				FileName:         fmt.Sprintf("README.md"),
				Parse:            false,
			},
			{
				TemplateFileName: "docker-compose.yml.tmpl",
				FileName:         fmt.Sprintf("docker-compose.yml"),
				Parse:            false,
			},
			{
				TemplateFileName: "Dockerfile.tmpl",
				FileName:         fmt.Sprintf("Dockerfile"),
				Parse:            false,
			},
			{
				TemplateFileName: "sqlboiler.toml.tmpl",
				FileName:         fmt.Sprintf("sqlboiler.toml"),
				Parse:            true,
			},
		}

		a.SetStructure(structure)

		if a.IsDirectoryExists(a.Project.Name) {
			log.Fatalf("chosen directory name already exists: %s", a.Project.Name)
		}

		syscall.Umask(0)
		err := os.Mkdir(a.Project.Name, 0755)
		if err != nil {
			a.Fatal(errors.Wrap(err, "error creating directory: %s"), a.Project.Name)

		}
		err = os.Chdir(a.Project.Name)
		if err != nil {
			a.Fatal(errors.Wrap(err, "error going into directory: %s"), a.Project.Name)
		}

		projectNamePath := fmt.Sprintf("cmd/%s", a.Project.Name)

		directories = append(directories, projectNamePath)
		err = a.CreateDirectories(directories)
		if err != nil {
			a.Fatal(err, a.Project.Name)
		}
		err = a.InitGoMod()
		if err != nil {
			a.Fatal(err)
		}

		err = a.CreateFiles()
		if err != nil {
			a.Fatal(err)
		}

		fmt.Printf(InfoColor, "...done.\n")
	},
}

func checkVersion() {
	version := runtime.Version()
	res := semver.Compare(version, "1.13")
	if res < 0 {
		log.Fatal(errors.New("go version must be > 1.13 for go modules support"))
	}
}
