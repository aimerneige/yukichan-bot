package config

var Conf *Config

type CommonConfig struct {
	WSServer      string   `yaml:"ws_server"`
	WSToken       string   `yaml:"ws_token"`
	NickName      []string `yaml:"nickname"`
	CommandPrefix string   `yaml:"command_prefix"`
	SuperUsers    []int64  `yaml:"super_users"`
}

type Config struct {
	Common   CommonConfig      `yaml:"common"`
	ConfFile map[string]string `yaml:"config"`
	DataFile map[string]string `yaml:"data"`
	TempDir  map[string]string `yaml:"temp"`
}
