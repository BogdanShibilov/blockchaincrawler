package config

type Config struct {
	GrpcServer GrpcServer `yaml:"grpcServer"`
	Auth       Auth       `yaml:"auth"`
	Transport  Transport  `yaml:"transport"`
	Kafka      Kafka      `yaml:"kafka"`
}

type GrpcServer struct {
	Port string `yaml:"port"`
}

type Auth struct {
	JwtSecretKey string `yaml:"jwtSecretKey"`
}

type Transport struct {
	UserGrpc UserGrpcTransport `yaml:"userGrpc"`
}

type UserGrpcTransport struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Kafka struct {
	Brokers  []string `yaml:"brokers"`
	Producer Producer `yaml:"producer"`
	Consumer Consumer `yaml:"consumer"`
}

type Producer struct {
	Topic string `yaml:"topic"`
}

type Consumer struct {
	Topics []string `yaml:"topics"`
}
