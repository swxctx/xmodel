package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Config db config
type Config struct {
	Database string `yaml:"database" json:"database"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	// the maximum number of connections in the idle connection pool.
	//
	// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns
	// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit
	//
	// If n <= 0, no idle connections are retained.
	MaxIdleConns int `yaml:"max_idle_conns" json:"max_idle_conns"`
	// the maximum number of open connections to the database.
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
	// MaxIdleConns, then MaxIdleConns will be reduced to match the new
	// MaxOpenConns limit
	//
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxOpenConns int `yaml:"max_open_conns" json:"max_open_conns"`
	// maximum amount of second a connection may be reused.
	// If d <= 0, connections are reused forever.
	ConnMaxLifetime int64 `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`

	// NoCache whether to disable cache
	NoCache bool `yaml:"no_cache" json:"no_cache"`
}

// NewConfig creates a default config.
func NewConfig() *Config {
	return &Config{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Database: "test",
	}
}

// Source returns the mysql connection string.
func (cfg *Config) Source() string {
	pwd := cfg.Password
	if pwd != "" {
		pwd = ":" + pwd
	}
	port := cfg.Port
	if port == 0 {
		port = 3306
	}
	return fmt.Sprintf("%s%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&interpolateParams=true", cfg.Username, pwd, cfg.Host, port, cfg.Database)
}
