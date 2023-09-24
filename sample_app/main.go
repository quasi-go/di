package main

import (
	"fmt"
	"github.com/quasi-go/di"
	"github.com/quasi-go/di/sample_app/config"
	"github.com/quasi-go/di/sample_app/services"
	"github.com/quasi-go/di/sample_app/sql"
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

	// We can already construct instances of `services.ServiceA` and `services.ServiceB`, even without explicit bindings.

	serviceA := di.Instance[services.ServiceA]()
	serviceB := di.Instance[services.ServiceB]()

	serviceA.DoSomethingInTheDatabase()
	serviceB.DoSomethingElseInTheDatabase()

	// Let's demonstrate binding concrete implementations to interfaces.

	di.BindType[services.IConfigToString, services.AppConfigToString]()
	//di.BindType[services.IConfigToString, services.DBConfigToString]() // uncomment this line to change what is printed below

	configReader := di.Instance[services.ConfigReader]()

	fmt.Println("This will print the bound IConfigToString :", configReader.Config.ToString())
}
