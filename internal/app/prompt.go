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

	prompt := promptui.Prompt{
		Label: "(optional) Type in URL without https://. Leave empty to use project name",
	}

	input, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	if input == "" {
		a.Project.ModuleName = a.Project.Name
		return
	}
	a.Project.ModuleName = input
}

func (a *App) inputDBInformation() {
	a.inputDBType()
	a.inputDriver()

	prompt := promptui.Prompt{
		Label: "Type in database address (default: 0.0.0.0)",
	}
	input, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	if input == "" {
		a.Project.Host = "0.0.0.0"
	} else {
		a.Project.Host = input
	}

	prompt = promptui.Prompt{
		Label: "Type in database name",
	}
	input, err = prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	a.Project.DBName = input

	prompt = promptui.Prompt{
		Label: "Type in database username",
	}
	input, err = prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	a.Project.Username = input

	prompt = promptui.Prompt{
		Label: "Type in database password",
	}
	input, err = prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	a.Project.Password = input

	promptSelect := promptui.Select{
		Label: "Enable SSL Mode? (disable)",
		Items: []string{"disable", "enable"},
	}
	_, result, err := promptSelect.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	if result == "" {
		a.Project.SSLMode = "disable"
	}
	a.Project.SSLMode = result
}

func (a *App) inputDBType() {
	prompt := promptui.Select{
		Label: "Select Database type (postgres",
		Items: []string{"postgres", "mariadb"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if result == "" || result == "postgres"{
		a.Project.Type = "postgres"
		a.Project.Port = 5432
	} else {
		a.Project.Type = result
		a.Project.Port = 3306
	}
}

func (a *App) inputDriver() {
	var defaultDriver string
	var driverOptions []string
	if a.Project.Type == "postgres" {
		defaultDriver = "sqlx/pq"
		driverOptions = []string{"sqlx/pq", "sqlboiler"}
		//driverOptions = []string{"sqlx/pq", "sqlboiler", "pgx"}
	} else {
		defaultDriver = "sqlx"
		driverOptions = []string{"sqlx", "database/sql"}
		//driverOptions = []string{"sqlx", "go-sql-driver/mysql", "database/sql"}
	}

	prompt := promptui.Select{
		Label: fmt.Sprintf("Select Database driver, (%s)", defaultDriver),
		Items: driverOptions,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	a.Project.Driver = result
}

func (a *App) scaffoldAuthentication() {
	prompt := promptui.Select{
		Label: "Scaffold authentication (no)",
		Items: []string{"no", "yes"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if result == "no" {
		a.Project.ScaffoldAuthentication = false
	}
	a.Project.ScaffoldAuthentication = true
}

func (a *App) printChoices() {
	fmt.Println("\nYou chose")
	fmt.Printf("Module Name            : %s\n", a.Project.ModuleName)
	fmt.Printf("Database Name          : %s\n", a.Project.DBName)
	fmt.Printf("Database Type          : %s\n", a.Project.Type)
	fmt.Printf("Database Driver        : %s\n", a.Project.Driver)
	fmt.Printf("Database Username      : %s\n", a.Project.Username)
	fmt.Printf("Database Password      : %s\n", a.Project.Password)
	fmt.Printf("Database SSL Mode      : %s\n", a.Project.SSLMode)
	fmt.Printf("Scaffold Authentication: %t\n", a.Project.ScaffoldAuthentication)
}
