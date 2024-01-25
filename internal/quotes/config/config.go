package config

import (
	"fmt"
	"time"
)

type Config struct {
	Server         Server  `mapstructure:"server"`
	Db             Db      `mapstructure:"db"`
	DaemonSettings Daemons `mapstructure:"daemons"`
}

type Server struct {
	Port int `mapstructure:"port"`
}

type Daemons struct {
	TaskRefresher DaemonSettings `mapstructure:"task_refresher"`
}
type DaemonSettings struct {
	BatchSize   int           `mapstructure:"batch_size"`
	BatchSleep  time.Duration `mapstructure:"sleep_batch"`
	NoWorkSleep time.Duration `mapstructure:"sleep_no_rows"`
}

type Db struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

func (db Db) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, db.User, db.Password, db.Name)
}
