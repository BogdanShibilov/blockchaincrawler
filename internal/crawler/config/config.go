package config

type Config struct {
	NodeUrl   NodeUrl   `yaml:"nodeUrl"`
	Transport Transport `yaml:"transport"`
}

type NodeUrl struct {
	Protocol string `yaml:"protocol"`
	Hostname string `yaml:"hostname"`
}

type Transport struct {
	BlockInfoTransport BlockInfoTransport `yaml:"blockInfoTransport"`
}

type BlockInfoTransport struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
