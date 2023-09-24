package main

import (
	"fmt"
	"github.com/quasi-go/di"
	"github.com/quasi-go/di/sample_app/config"
	s "github.com/quasi-go/di/sample_app/services"
	"github.com/quasi-go/di/sample_app/sql" // this is a stub of "database/sql"
	"log"
	"os"
)

func main() {
	// Set logger
	logger := log.New(os.Stdout, "DI: ", 0)
	di.SetLogger(logger)

	// This would usually involve reading from environment variable instead of hard-coded strings.

	di.BindInstance(&config.AppConfig{
		VarA: "This will be sent as an argument to the query from ServiceA",
		VarB: "This will be sent as an argument to the query from ServiceB",
	})

	di.BindInstance(&config.DBConfig{
		Driver:   "driver",
		Host:     "host",
		Port:     "123",
		Username: "username",
		Password: "password",
		DBName:   "dbname",
	})

	// Now we bind a provider that will use the DBConfig we bound to above.

	di.BindProvider(func(config config.DBConfig) (*sql.DB, error) {
		db, err := sql.Open(
			config.Driver,
			config.Driver+"://"+config.Username+":"+config.Password+
				"@"+config.Host+":"+config.Port+"/"+config.DBName,
		)

		if err != nil {
			return nil, err
		}

		err = db.Ping()

		if err != nil {
			return nil, err
		}

		return db, nil
	})

	// We can already construct instances of `ServiceA` and `ServiceB`, even without explicit bindings.

	serviceA := di.Instance[s.ServiceA]()
	serviceB := di.Instance[s.ServiceB]()

	serviceA.DoSomethingInTheDatabase()
	serviceB.DoSomethingElseInTheDatabase()

	// Let's demonstrate binding concrete implementations to interfaces.

	di.BindType[s.IConfigToString, s.AppConfigToString]()
	//di.BindType[s.IConfigToString, s.DBConfigToString]() // uncomment this line to change what is printed below

	configReader := di.Instance[s.ConfigReader]()

	fmt.Println("This will print the bound IConfigToString :", configReader.Config.ToString())
}
