package config

type AppConfig struct {
	VarA string
	VarB string
}

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}
