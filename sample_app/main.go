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

	// Below we bind our config instances that will be used to populated dependencies
	// other types we utilize in our example. This would usually involve reading from
	// environment variable instead of hard-coded strings.

	di.BindInstance(&config.AppConfig{
		VarA: "This will be sent as an argument to the query from ServiceA",
		VarB: "This will be sent as an argument to the query from ServiceA",
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
	// We `db.Open()` and then `db.Ping()` a connection, checking for errors along the way.
	// This binding will provide the returned type `sql.Db` for other types.

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

	// We can already construct instances of `services.ServiceA` and `services.ServiceB`
	// even without explicit bindings.

	serviceA := di.Instance[services.ServiceA]()
	serviceB := di.Instance[services.ServiceB]()

	serviceA.DoSomethingInTheDatabase()
	serviceB.DoSomethingElseInTheDatabase()

	// Let's demonstrate binding concrete implementations to interfaces.

	di.BindType[services.IConfigToString, services.AppConfigToString]()
	configToString := di.Impl[services.IConfigToString]() // Note that we use `Impl[T]()` here instead of `Instance[T]()`

	fmt.Println("This will print AppConfig:", configToString.ToString())

	// We can override our previous binding to `services.IConfigToString`.

	di.BindType[services.IConfigToString, services.DBConfigToString]()
	configToString = di.Impl[services.IConfigToString]() // Note that we use `Impl[T]()` here instead of `Instance[T]()`

	fmt.Println("This will print DBConfig:", configToString.ToString())

	di.GetContainer().SetLogger(nil)

	// Finally, we can demonstrate implicit construction of structs with interface dependencies.
	// The `services.IConfigToString` dependency will be satisfied by our most recent binding to `services.DBConfigToString`.

	// Note that even though it had an interface dependency,
	// `services.ConfigReader` is a struct to `Instance[T]()` is used.

	configReader := di.Instance[services.ConfigReader]()

	fmt.Println("This will print DBConfig again:", configReader.Config.ToString())
}
