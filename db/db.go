package db

// Config is a configuration instance for connecting DB
type Config struct {
	Username string
	Password string
	Host     string
	Port     int
}

type PersistentService interface {
	GenerateConnectLink() (link string, err error)
	InitTable() error
	CommitId() error
	ConsumeId() (id int64, err error)
	BatchCommitId() error
	BatchConsumeId() (idList []int64, err error)
	CountUnusedId() (num int, err error)
}
