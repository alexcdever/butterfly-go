package db

type MysqlConfig struct {
	Config
	encoding string
}

func (config *MysqlConfig) Connect() {

}
