package agent

import (
	"encoding/json"
	"errors"
	"github.com/percona/cloud-protocol/proto"
	"io/ioutil"
	"log"
	"os"
)

// Defaults
const (
	API_HOSTNAME = "cloud-api.percona.com"
	CONFIG_FILE  = "/etc/percona/agent.conf"
	LOG_FILE     = "/var/log/percona/agent.log"
	LOG_LEVEL    = "info"
	DATA_DIR     = "/var/spool/percona/agent"
)

type Config struct {
	ApiHostname   string
	ApiKey        string
	AgentUuid     string
	PidFile       string
	LogFile       string
	LogLevel      string
	DataDir       string
	Links         map[string]string
	Enable        []string
	Disable       []string
}

// Load config from JSON file.
func LoadConfig(file string) *Config {
	config := new(Config)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalln(err)
		}
	} else {
		if err = json.Unmarshal(data, config); err != nil {
			log.Fatalln(err)
		}
	}
	return config
}

// Apply current config, i.e. overwrite this config with current config.
func (c *Config) Apply(cur *Config) error {
	if cur.ApiHostname != "" {
		c.ApiHostname = cur.ApiHostname
	}
	c.ApiKey = cur.ApiKey
	c.AgentUuid = cur.AgentUuid
	c.PidFile = cur.PidFile
	if cur.LogFile != "" {
		c.LogFile = cur.LogFile
	}
	if cur.LogLevel != "" {
		_, ok := proto.LogLevels[cur.LogLevel]
		if !ok {
			return errors.New("Invalid log level: " + cur.LogLevel)
		}
		c.LogLevel = cur.LogLevel
	}
	c.DataDir = cur.DataDir
	c.Enable = cur.Enable
	c.Disable = cur.Disable
	return nil
}

func (c *Config) Enabled(option string) bool {
	for _, o := range c.Enable {
		if o == option {
			return true
		}
	}
	return false
}

func (c *Config) Disabled(option string) bool {
	for _, o := range c.Disable {
		if o == option {
			return true
		}
	}
	return false
}
