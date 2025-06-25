package initialize

func Run() {
	LoadConfig()
	InitPostgresDB()
	InitRedis()
	InitLogger()
}
