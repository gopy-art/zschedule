package configs

type CommandLineConfig struct {
	Command  string `yaml:"COMMAND" json:"command"`
	Interval int    `yaml:"INTERVAL" json:"interval"`
	Limit    int    `yaml:"LIMIT" json:"limit"`
}
