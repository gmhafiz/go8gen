package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/friendsofgo/errors"
)

func (p *Project) InjectCode() error {
	const serverFileName = "internal/server/server.go"
	const injectImport = "// inject:import"
	//const injectApp = "//inject:app"
	//const injectUseCase = "// inject:usecase"
	const injectHandler = "// inject:handler"
	//p.Domain = strings.Title(p.DomainLowerCase)
	importTmpl1 := fmt.Sprintf(`%sHTTP "%s/internal/domain/%s/handler/http"`, p.DomainLowerCase, p.Name, p.DomainLowerCase)
	//importTmpl2 := fmt.Sprintf(`%sPostgres "%s/internal/domain/%s/repository/postgres"`, p.DomainLowerCase, p.Name, p.DomainLowerCase)
	//importTmpl3 := fmt.Sprintf(`%sUseCase "%s/internal/domain/%s/usecase"`, p.DomainLowerCase, p.Name, p.DomainLowerCase)
	//appTmpl := fmt.Sprintf(`%sUC *%sUseCase.%sUseCase`, p.DomainLowerCase, p.DomainLowerCase, p.Domain)
	//usecaseTmpl := fmt.Sprintf(`%sUC: %sUseCase.New%sUseCase(%sPostgres.New%sRepository(db)),`, p.DomainLowerCase, p.DomainLowerCase, p.Domain, p.DomainLowerCase, p.Domain)
	handlerTmpl := fmt.Sprintf(`%sHTTP.RegisterHTTPEndPoints(router, a.%sUC)`, p.DomainLowerCase, p.DomainLowerCase)

	serverContent, err := ioutil.ReadFile(serverFileName)
	if err != nil {
		return errors.Wrapf(err, "error reading file: %s", serverFileName)
	}

	var newFile []string
	temp := strings.Split(string(serverContent), "\n")
	for _, line := range temp {
		newFile = append(newFile, line)
		stripped := strings.Trim(line, "\t")
		stripped = strings.Trim(stripped, "\n")
		if stripped == injectImport {
			newFile = append(newFile, importTmpl1)
			//newFile = append(newFile, importTmpl2)
			//newFile = append(newFile, importTmpl3)
		}
		//if stripped == injectApp {
		//	newFile = append(newFile, appTmpl)
		//}
		//if stripped == injectUseCase {
		//	newFile = append(newFile, usecaseTmpl)
		//}
		if stripped == injectHandler {
			newFile = append(newFile, handlerTmpl)
		}
	}

	fCreate, err := os.Create(serverFileName)
	if err != nil {
		return errors.Wrapf(err, "error opening file: %s", serverFileName)
	}
	_, err = fCreate.WriteString(strings.Join(newFile, "\n"))
	if err != nil {
		return errors.Wrapf(err, "error writing file: %s", serverFileName)
	}
	return nil
}

func (p *Project) InjectImportDomainHandlerCode() error {
	const initDomainsFile = "internal/server/initDomains.go"
	const injectImport = "import ("
	const injectInitDomains = "func (s *Server) InitDomains() {"

	initDomainsTmpl := fmt.Sprintf(`    s.init%s()`, p.Domain)

	var importTmpl1 string
	if p.ScaffoldUseCase {
		importTmpl1 = fmt.Sprintf(`	%sHandler "%s/internal/domain/%s/handler/http"
	%sRepo "%s/internal/domain/%s/repository/database"
	%sUseCase "%s/internal/domain/%s/usecase"`,
			p.DomainLowerCase, p.ModuleName, p.DomainLowerCase,
			p.DomainLowerCase, p.ModuleName, p.DomainLowerCase,
			p.DomainLowerCase, p.ModuleName, p.DomainLowerCase,
		)
	} else {
		importTmpl1 = fmt.Sprintf(`    new%sUseCase := %sUseCase.New$sUseCase()
%sHandler "%s/internal/domain/%s/handler/http"`, p.Domain, p.DomainLowerCase, p.DomainLowerCase,
			p.Name, p.DomainLowerCase)
	}

	var initHandlerTmpl string
	if p.ScaffoldUseCase {
		initHandlerTmpl = fmt.Sprintf(`
func (s *Server) init%s() {
	new%sRepo := %sRepo.New(s.DB())
	new%sUseCase := %sUseCase.New(new%sRepo)
	%sHandler.RegisterHTTPEndPoints(s.router, s.validator, new%sUseCase)
}`, p.Domain,
			p.Domain, p.DomainLowerCase,
			p.Domain, p.DomainLowerCase, p.Domain,
			p.DomainLowerCase, p.Domain)
	} else {
		initHandlerTmpl = fmt.Sprintf(`func (s *Server) init%s() {
	%sHandler.RegisterHTTPEndPoints(s.router)
}`, p.Domain, strings.ToLower(p.Domain))
	}

	fileContent, err := ioutil.ReadFile(initDomainsFile)
	if err != nil {
		return errors.Wrapf(err, "error reading file: %s", initDomainsFile)
	}

	var newFile []string
	temp := strings.Split(string(fileContent), "\n")
	for _, line := range temp {
		newFile = append(newFile, line)
		stripped := strings.Trim(line, "\t")
		stripped = strings.Trim(stripped, "\n")
		if stripped == injectImport {
			newFile = append(newFile, importTmpl1)
		}
		if stripped == injectInitDomains {
			newFile = append(newFile, initDomainsTmpl)
		}
	}
	newFile = append(newFile, initHandlerTmpl)

	fCreate, err := os.Create(initDomainsFile)
	if err != nil {
		return errors.Wrapf(err, "error opening file: %s", initDomainsFile)
	}
	_, err = fCreate.WriteString(strings.Join(newFile, "\n"))
	if err != nil {
		return errors.Wrapf(err, "error writing file: %s", initDomainsFile)
	}
	return nil
}
