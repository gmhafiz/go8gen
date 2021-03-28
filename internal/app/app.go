package app

import (
	"embed"
	"os/exec"

	"github.com/friendsofgo/errors"
)

//go:embed tmpl/*
var static embed.FS

type Project struct {
	Name            string
	ModuleName      string
	Domain          string
	DomainLowerCase string

	ScaffoldAuthentication bool
	ScaffoldUseCase        bool
	ScaffoldRepository     bool

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
		Project:   Project{},
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
