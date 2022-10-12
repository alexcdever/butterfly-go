package db

type MysqlConfig struct {
	Config
	encoding string
}

func (config *MysqlConfig) GenerateConnectLink() (link string, err error) {
	return link, err
}
