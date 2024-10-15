package config

import (
	"fmt"
	"os"
)

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	Database Database
}

// DbConnURL returns string URL, which may be used for connect to postgres database.
func (c *Database) ConnURL() string {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
	return url
}

func NewFrom(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		// logrus.WithError(err).WithField("path", path).Fatal("can't read config file")
	}
	defer file.Close()
	var config = Config{}
	// Init new YAML decode
	// d := yaml.NewDecoder(file)
	// Start YAML decoding from file
	// if err := d.Decode(&config); err != nil {
	// logrus.WithError(err).Fatal("can't decode config file")
	// }
	// if err := binding.Validator.ValidateStruct(config); err != nil {
	// logrus.WithError(err).Fatal("config validation failed")
	// }
	// level, err := logrus.ParseLevel(config.LogLevel)
	// if err != nil {
	// 	logrus.Fatal("invalid 'logLevel' parameter in configuration. Available values: ", logrus.AllLevels)
	// }
	// logrus.SetLevel(level)
	// logrus.SetReportCaller(true) // adds line number to log message
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	return config
}
