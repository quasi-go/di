package services

import (
	"encoding/json"
	"github.com/quasi-go/di/sample_app/config"
	"github.com/quasi-go/di/sample_app/sql"
)

// Below we have a couple sample services ServiceA & Service B that consume our
// DB handle and our AppConfig. Note that these dependencies can be filled by
// reference or by value based on whether they are defined as pointers.

type ServiceA struct {
	DB     *sql.DB
	Config config.AppConfig
}

func (a *ServiceA) DoSomethingInTheDatabase() {
	a.DB.Query("Query from ServiceA", a.Config.VarA)
}

type ServiceB struct {
	DB     *sql.DB
	Config config.AppConfig
}

func (b *ServiceB) DoSomethingElseInTheDatabase() {
	b.DB.Query("Query from ServiceB", b.Config.VarB)
}

// Below we have a conjured use case of converting our two config struct
// types, AppConfig and DBConfig to string. We create two structs that
// implement the interface IConfigToString to demonstrate how to by a
// concrete implementation to an interface.

type IConfigToString interface {
	ToString() string
}

type AppConfigToString struct {
	AppConfig config.AppConfig
}

func (c *AppConfigToString) ToString() string {
	encoded, _ := json.Marshal(c.AppConfig)
	return string(encoded)
}

type DBConfigToString struct {
	DBConfig config.DBConfig
}

func (c *DBConfigToString) ToString() string {
	encoded, _ := json.Marshal(c.DBConfig)
	return string(encoded)
}

type ConfigReader struct {
	Config IConfigToString
}
