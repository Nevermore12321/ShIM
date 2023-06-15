package confs

import "time"

type MySQL struct {
	Read  ConnOptions `yaml:"read"`
	Write ConnOptions `yaml:"write"`
	Base  BaseOptions `yaml:"base"`
}

type ConnOptions struct {
	Hosts    []string `yaml:"host"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Name     string   `yaml:"name"`
}

type BaseOptions struct {
	MaxOpenConn     int           `yaml:"maxOpenConn"`
	MaxIdleConn     int           `yaml:"maxIdleConn"`
	ConnMaxLifeTime time.Duration `yaml:"connMaxLifeTime"`
	DatabaseName    string        `yaml:"dbMysqlDatabaseName"`
	DBTableName     string        `yaml:"DBTableName"`
	DBMsgTableNum   int           `yaml:"dbMsgTableNum"`
	LogLevel        int           `yaml:"logLevel"`
	SlowThreshold   int           `yaml:"slowThreshold"`
}
