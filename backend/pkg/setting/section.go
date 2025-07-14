package setting

import "time"

type Server struct {
	Port int `mapstructure:"port"`
}

type Redis struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Password        string `mapstructure:"password"`
	DB              int    `mapstructure:"db"`
	PoolSize        int    `mapstructure:"poolSize"`
	MinIdleConns    int    `mapstructure:"minIdleConns"`
	DialTimeout     int    `mapstructure:"dialTimeout"`
	ReadTimeout     int    `mapstructure:"readTimeout"`
	WriteTimeout    int    `mapstructure:"writeTimeout"`
	ConnMaxIdleTime int    `mapstructure:"connMaxIdleTime"`
}

type Mysql struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DbName          string `mapstructure:"dbName"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
	ConnMaxIdleTime int    `mapstructure:"connMaxIdleTime"`
}

type JWT struct {
	Secret         string        `mapstructure:"secret"`
	ExpirationTime time.Duration `mapstructure:"expiration_time"`
	EncryptionKey  string        `mapstructure:"encryption_key"`
}

type Cookie struct {
	Domain   string `mapstructure:"domain"`
	Secure   bool   `mapstructure:"secure"`
	HttpOnly bool   `mapstructure:"http_only"`
	SameSite string `mapstructure:"same_site"`
}

type RabbitMQ struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	IamExchange string `mapstructure:"iam_exchange"`
}

type Config struct {
	Server   Server   `mapstructure:"server"`
	Redis    Redis    `mapstructure:"redis"`
	Mysql    Mysql    `mapstructure:"mysql"`
	JWT      JWT      `mapstructure:"jwt"`
	Cookie   Cookie   `mapstructure:"cookie"`
	RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
}
