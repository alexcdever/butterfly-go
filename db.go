package main

// DbConfig is a configuration instance for connecting DB
type DbConfig struct {
	username string
	password string
	url      string
}

type MysqlConfig struct {
	DbConfig
	encoding string
}
type DbConfigService interface {
	Connect()
	InitDB()
	GenerateId()
	ConsumeId()
}

func (config *MysqlConfig) Connect() {

}
