package app

//go:generate binclude

import (
	"github.com/lu4p/binclude"
)

type Project struct {
	Name            string
	ModuleName      string
	Domain          string
	DomainLowerCase string

	ScaffoldAuthentication bool

	Type     string
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
}


type Structure struct {
	TemplateFileName string
	FileName         string
	Parse            bool
}

type App struct {
	Project   Project
	Structure []Structure
	TemplatePath string
}

func New() *App {
	return &App{
		Project:   Project{},
		Structure: []Structure{},
		TemplatePath: binclude.Include("../../tmpl"),
	}
}

func (a *App) SetStructure(structure []Structure) {
	a.Structure = structure
}

func (a *App) SetProject(p Project) {
	a.Project = p
}


