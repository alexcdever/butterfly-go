package db

type PostgreSqlConfig struct {
	Config
	Postgres struct {
		Database string
		SslMode  bool
		TimeZone string
	}
}

func (config *PostgreSqlConfig) GenerateConnectLink() (link string, err error) {
	return link, err
}
