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

type Config struct {
	Server Server `mapstructure:"server"`
	Redis  Redis  `mapstructure:"redis"`
	Mysql  Mysql  `mapstructure:"mysql"`
}
