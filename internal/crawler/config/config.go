package config

type Config struct {
	NodeUrl NodeUrl `yaml:"nodeUrl"`
}

type NodeUrl struct {
	Protocol string `yaml:"protocol"`
	Hostname string `yaml:"hostname"`
}
