package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"time"
)

// MySQLOpts provides config for mysql.
// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
// https://github.com/go-sql-driver/mysql#dsn-data-source-name
type MySQLOpts struct {
	Username string `json:"username,omitempty" mapstructure:"username"`
	Password string `json:"password,omitempty" mapstructure:"password"`
	Host     string `json:"host,omitempty" mapstructure:"host"`
	Port     int    `json:"port,omitempty" mapstructure:"port"`
	DbName   string `json:"db_name,omitempty" mapstructure:"db_name"`
	Pool     *Pool  `json:"pool" mapstructure:"pool"`
	LogLevel int    `json:"log_level,omitempty" mapstructure:"log_level"`
}

type Pool struct {
	MaxIdleConnections    int           `json:"max_idle_connections,omitempty" mapstructure:"max_idle_connections"`
	MaxOpenConnections    int           `json:"max_open_connections,omitempty" mapstructure:"max_open_connections"`
	MaxConnectionLifeTime time.Duration `json:"max_connection_life_time,omitempty" mapstructure:"max_connection_life_time"`
}

func NewMySQLOpts() *MySQLOpts {
	return &MySQLOpts{
		Username: "",
		Password: "",
		Host:     "127.0.0.1",
		Port:     3306,
		DbName:   "",
		Pool: &Pool{
			MaxIdleConnections:    100,
			MaxOpenConnections:    100,
			MaxConnectionLifeTime: 10 * time.Second,
		},
		LogLevel: 0,
	}
}

func (o *MySQLOpts) DSN() string {
	return fmt.Sprintf("%s:%d@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		o.Username,
		o.Password,
		o.Host,
		o.Port,
		o.DbName)
}

func (o *MySQLOpts) Validate() []error {
	var err []error

	return err
}

func (o *MySQLOpts) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "mysql.host", o.Host, ""+
		"MySQL service host address. If left blank, the following related mysql options will be ignored.")

	fs.IntVar(&o.Port, "mysql.host", o.Port, ""+
		"MySQL service host port. If left blank, the following related mysql options will be ignored.")

	fs.StringVar(&o.Username, "mysql.username", o.Username, ""+
		"Username for access to mysql service.")

	fs.StringVar(&o.Password, "mysql.password", o.Password, ""+
		"Password for access to mysql, should be used pair with password.")

	fs.StringVar(&o.DbName, "mysql.database", o.DbName, ""+
		"Database name for the server to use.")

	fs.IntVar(&o.Pool.MaxIdleConnections, "mysql.pool.max-idle-connections", o.Pool.MaxOpenConnections, ""+
		"Maximum idle connections allowed to connect to mysql.")

	fs.IntVar(&o.Pool.MaxOpenConnections, "mysql.pool.max-open-connections", o.Pool.MaxOpenConnections, ""+
		"Maximum open connections allowed to connect to mysql.")

	fs.DurationVar(&o.Pool.MaxConnectionLifeTime, "mysql.pool.max-connection-life-time", o.Pool.MaxConnectionLifeTime, ""+
		"Maximum connection life time allowed to connect to mysql.")
}
