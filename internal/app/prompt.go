package app

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func (a *App) Prompt() {
	a.inputModuleName()
	a.inputDBInformation()
	a.scaffoldAuthentication()
	a.printChoices()
}

func (a *App) inputModuleName() {
	if a.Project.ModuleName != "" {
		return
	}

	input := prompt(fmt.Sprintf("(optional) Set Project URL without https:// or leave empty to use project name (%s)", a.Project.Name))
	if input == "" {
		a.Project.ModuleName = a.Project.Name
		return
	}
	a.Project.ModuleName = input
}

func (a *App) inputDBInformation() {
	proceed := a.willInputDBType()
	if proceed {
		a.inputDBType()
	}

	proceed = a.willInputDBInformation()
	if !proceed {
		return
	}

	a.inputDriver()

	input := prompt("Type in database address (default: 0.0.0.0)")
	if input != "" {
		a.Project.Host = input
	}

	input = prompt("Type in database name")
	if input != "" {
		a.Project.DBName = input
	}

	input = prompt("Type in database username")
	if input != "" {
		a.Project.Username = input
	}

	input = prompt("Type in database password")
	if input != "" {
		a.Project.Password = input
	}

	input = promptSelect("Enable SSL Mode? (disable)", []string{"disable", "enable"})
	if input != "" {
		a.Project.SSLMode = input
	}
}

func (a *App) willInputDBType() bool {
	result := promptSelect("Choose Database Type", []string{"no", "yes"})
	if result == "no" {
		return false
	} else {
		return true
	}
}

func (a *App) willInputDBInformation() bool {
	result := promptSelect("Set Database Credentials?", []string{"no", "yes"})
	if result == "no" {
		fmt.Println("Remember to fill in .env file if you going to use a database!")
		return false
	} else {
		return true
	}
}

func (a *App) inputDBType() {
	result := promptSelect("Select Database type (postgres)", []string{"postgres", "mysql"})
	switch result {
	case "postgres":
		a.Project.Type = "postgres"
		a.Project.Port = 5432
		a.Project.SqlBoilerDriverName = "psql"
	case "mysql":
		a.Project.Type = result
		a.Project.Port = 3306
		a.Project.SqlBoilerDriverName = "mysql"
	case "mssql":
		a.Project.Type = result
		a.Project.Port = 1433
		a.Project.SqlBoilerDriverName = "mssql"
	}
}

func (a *App) inputDriver() {
	var defaultDriver string
	var driverOptions []string

	switch a.Project.Type {
	case "postgres":
		defaultDriver = "sqlx/pgx"
		driverOptions = []string{"sqlx/pgx", "ent"}
	case "mysql":
		defaultDriver = "sqlx"
		driverOptions = []string{"sqlx", "database/sql"}
	}

	a.Project.Driver = promptSelect(fmt.Sprintf("Select Database driver, (%s)", defaultDriver), driverOptions)
}

func (a *App) scaffoldAuthentication() {
	result := promptSelect("Scaffold authentication", []string{"no", "yes"})
	if result == "no" {
		a.Project.ScaffoldAuthentication = false
	}
	a.Project.ScaffoldAuthentication = true
}

func (a *App) printChoices() {
	fmt.Println("\nYou chose")
	fmt.Println("---------")
	fmt.Printf("\nProject Path         : %s\n", a.Project.Path)
	fmt.Printf("Module Name            : %s\n", a.Project.ModuleName)
	fmt.Printf("Database Name          : %s\n", a.Project.DBName)
	fmt.Printf("Database Type          : %s\n", a.Project.Type)
	fmt.Printf("Database Driver        : %s\n", a.Project.Driver)
	fmt.Printf("Database Username      : %s\n", a.Project.Username)
	fmt.Printf("Database Password      : %s\n", a.Project.Password)
	fmt.Printf("Database SSL Mode      : %s\n", a.Project.SSLMode)
	fmt.Printf("Scaffold Authentication: %t\n", a.Project.ScaffoldAuthentication)
}

func prompt(message string) string {
	prompt := promptui.Prompt{
		Label: message,
	}
	input, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}
	return input
}

func promptSelect(message string, options []string) string {
	promptSelect := promptui.Select{
		Label: message,
		Items: options,
	}
	_, result, err := promptSelect.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}
	return result
}
