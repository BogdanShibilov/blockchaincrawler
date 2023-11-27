package config

type Config struct {
	Server       Server       `yaml:"server"`
	ExternalNode ExternalNode `yaml:"externalNode"`
	Transport    Transport    `yaml:"transport"`
}

type Server struct {
	Port string `yaml:"port"`
}

type ExternalNode struct {
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
