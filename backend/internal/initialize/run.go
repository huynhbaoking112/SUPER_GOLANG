package initialize

func Run() {
	LoadConfig()
	InitMysql()
	InitRedis()
	InitLogger()
}
