package services

import (
	"fmt"
	"github.com/quasi-go/di/sample_app/config"
	"github.com/quasi-go/di/sample_app/sql"
	"testing"
)

func TestServices(t *testing.T) {
	db := &sql.DB{}
	appConfig := config.AppConfig{}
	dbConfig := config.DBConfig{}

	a := ServiceA{db, appConfig}
	b := ServiceB{db, appConfig}

	a.DoSomethingInTheDatabase()
	b.DoSomethingElseInTheDatabase()

	acts := AppConfigToString{appConfig}
	dcts := DBConfigToString{dbConfig}
	cr := ConfigReader{&acts}

	fmt.Println(acts.ToString())
	fmt.Println(dcts.ToString())
	fmt.Println(cr.Config.ToString())
}
