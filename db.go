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
type DbPersistentService interface {
	Connect() error
	InitDB() error
	GenerateId() error
	ConsumeId() (id int64, err error)
	CheckUnusedId() (num int, err error)
}

func (config *MysqlConfig) Connect() {

}
