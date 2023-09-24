package sql

import (
	"github.com/quasi-go/di/sample_app/config"
	"testing"
)

func TestDB(t *testing.T) {
	dbConfig := &config.DBConfig{
		Driver:   "driver",
		Host:     "host",
		Port:     "123",
		Username: "username",
		Password: "password",
		DBName:   "dbname",
	}

	db, _ := Open(
		dbConfig.Driver,
		dbConfig.Driver+"://"+dbConfig.Username+":"+dbConfig.Password+
			"@"+dbConfig.Host+":"+dbConfig.Port+"/"+dbConfig.DBName,
	)

	_ = db.Ping()
}
