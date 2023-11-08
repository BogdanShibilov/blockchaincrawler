package config

type Config struct {
	GrpcServer GrpcServer `yaml:"grpcServer"`
	Database   Database   `yaml:"database"`
}

type GrpcServer struct {
	Port string `yaml:"port"`
}

type Database struct {
	Main DbNode `yaml:"main"`
}

type DbNode struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SslMode  string `yaml:"sslMode"`
}
