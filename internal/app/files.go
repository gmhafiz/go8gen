package app

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"text/template"

	"github.com/friendsofgo/errors"
)

const (
	tmplPath = "tmpl"
	tmpPath  = "/tmp"
)

func (a *App) CreateDirectories(directories []string) error {
	for _, val := range directories {
		syscall.Umask(0)
		err := os.MkdirAll(val, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) IsDirectoryExists(projectName string) bool {
	_, err := os.Stat(projectName)
	if err != nil {
		return false
	}
	return true
}

func (a *App) Fatal(msg error, args ...string) {
	err := os.RemoveAll(a.Project.Name)
	if err != nil {
		log.Printf("error removing diectory: %s", a.Project.Name)
	}
	log.Fatalf(msg.Error(), args)
}

func (a *App) InitGoMod() error {
	cmd := exec.Command("go", "mod", "tidy")
	//cmd.Dir = a.Project.Path
	cmd.Path = a.Project.Path
	_, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "error running: go mod tidy")
	}
	return nil
}

func (a *App) CreateFiles() error {
	for _, val := range a.Structure {
		err := a.CreateFile(val.TemplateFileName, val.FileName, val.Parse)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) CreateFile(tmplFileName, fileName string, parse bool) error {
	file, err := os.Create(fileName)
	if err != nil {
		return errors.Wrapf(err, "error creating file: %s", fileName)
	}

	if !parse {
		return a.CopyFile(fileName, tmplFileName)
	} else {
		return a.ParseFile(file, tmplFileName)
	}
}

func (a *App) ParseFile(file *os.File, tmplFileName string) error {
	tmplContent, err := a.Static.ReadFile(filepath.Join(tmplPath, tmplFileName))
	if err != nil {
		return errors.Wrap(err, "error reading template file")
	}
	fileNameTail := filepath.Base(tmplFileName)
	f, err := os.Create(filepath.Join(tmpPath, fileNameTail))
	if err != nil {
		return errors.Wrap(err, "error creating file at /tmp folder")
	}
	_, err = f.Write(tmplContent)
	if err != nil {
		return errors.Wrap(err, "error writing temporary file")
	}

	tmpl, err := template.ParseFiles(filepath.Join(tmpPath, fileNameTail))
	if err != nil {
		return errors.Wrapf(err, "error parsing file: %s", fileNameTail)
	}
	return tmpl.Execute(file, a.Project)
}

func (a *App) CopyFile(fileName, tmplFileName string) error {
	tmplContent, err := a.Static.ReadFile(filepath.Join(tmplPath, tmplFileName))
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
