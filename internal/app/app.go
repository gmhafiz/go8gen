package app

import (
	"embed"
	"os/exec"

	"github.com/friendsofgo/errors"
)

//go:embed tmpl/*
var static embed.FS

type TaskFile struct {
	CliArgs string
	N       string
}

type Project struct {
	Path            string
	Address         string
	Name            string
	ModuleName      string
	Domain          string
	DomainLowerCase string

	Host string
	Port int

	ScaffoldAuthentication bool
	ScaffoldUseCase        bool
	ScaffoldRepository     bool

	Database

	TaskFile

	ExpandedID bool
}

type Database struct {
	Type     string
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string

	SqlBoilerDriverName string
}

type Structure struct {
	TemplateFileName string
	FileName         string
	Parse            bool
}

type App struct {
	Project   Project
	Structure []Structure
	Static    embed.FS
}

func New() *App {
	return &App{
		Project: Project{
			Database: Database{
				Type:     "postgres",
				Driver:   "postgres",
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "",
				DBName:   "postgres",
			},
			TaskFile: TaskFile{
				CliArgs: "{{.CLI_ARGS}}",
				N:       "{{.n}}",
			},
		},
		Structure: []Structure{},
		Static:    static,
	}
}

func (a *App) SetStructure(structure []Structure) {
	a.Structure = structure
}

func (a *App) SetProject(p *Project) {
	a.Project = *p
}

func (a *App) GoModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed with %s\n", err)
	}

	cmd = exec.Command("go", "fmt", "./...")
	_ = cmd.Run() // ignoring output because we do not care files that have been formatted

	return nil
}
