package setting

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

type Config struct {
	Server Server `mapstructure:"server"`
	Redis  Redis  `mapstructure:"redis"`
}
