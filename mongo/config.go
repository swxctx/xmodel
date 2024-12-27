package mongo

import (
	"time"

	"gopkg.in/mgo.v2"
)

// Config mongodb setting
type Config struct {
	Addrs    []string      `yaml:"addrs" json:"addrs"`
	Timeout  time.Duration `yaml:"timeout" json:"timeout"`
	Database string        `yaml:"database" json:"database"`
	// Username and Password inform the credentials for the initial authentication
	// done on the database defined by the Source field. See Session.Login.
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"passward" json:"password"`
	// PoolLimit defines the per-server socket pool limit. Defaults to 4096.
	// See Session.SetPoolLimit for details.
	PoolLimit int `yaml:"pool_limit" json:"pool_limit"`
	// NoCache whether to disable cache
	NoCache bool `yaml:"no_cache" json:"no_cache"`
}

// NewConfig creates a default config.
func NewConfig() *Config {
	return &Config{
		Addrs:     []string{"127.0.0.1:27017"},
		Timeout:   10,
		PoolLimit: 256,
		Username:  "root",
		Password:  "",
		Database:  "test",
	}
}

// Set config
func (mgoConfig *Config) Source() *mgo.DialInfo {
	dialInfo := &mgo.DialInfo{
		Addrs:     mgoConfig.Addrs,
		Username:  mgoConfig.Username,
		Password:  mgoConfig.Password,
		Database:  mgoConfig.Database,
		Timeout:   mgoConfig.Timeout,
		PoolLimit: mgoConfig.PoolLimit,
	}

	return dialInfo
}
