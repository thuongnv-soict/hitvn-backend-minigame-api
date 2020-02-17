package main

import (
	"fmt"
	"g-tech.com/infrastructure/broker"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/infrastructure/repository"
	"g-tech.com/infrastructure/util"
	"g-tech.com/module/lottery"
	"github.com/spf13/viper"
	"os"
	"time"
)

func init(){
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}
}
func main() {
	/********************************************************************/
	/* CONFIGURE LOG													*/
	/********************************************************************/
	logPath 	:= viper.GetString(`Log.Path`)
	logPrefix 	:= viper.GetString(`Log.PrefixLottery`)
	logger.NewLogger(logPath, logPrefix)

	timeout := time.Duration(viper.GetInt("Context.Timeout")) * time.Second
	/********************************************************************/
	/* CONFIGURE RabbitMQ												*/
	/********************************************************************/
	// Load RabbitMQ configuration
	rabbitHost				:= viper.GetString(`RabbitMQ.Host`)
	rabbitPort				:= viper.GetInt(`RabbitMQ.Port`)
	rabbitUserName			:= viper.GetString(`RabbitMQ.UserName`)
	rabbitPassword			:= viper.GetString(`RabbitMQ.Password`)

	// Open a RabbitMQ
	rbConnection := broker.Connect(rabbitHost, rabbitPort, rabbitUserName, rabbitPassword)
	if rbConnection == nil {
		logger.Panic("Failed to connect to RabbitMQ")
		os.Exit(1)
	}
	defer rbConnection.Close()

	// Creates a RabbitMQ channel
	rbChannel, err := rbConnection.Channel()
	if err != nil {
		logger.Error(err.Error())
	}
	defer rbChannel.Close()

	// Set RabbitMQ QoS
	err = rbChannel.Qos(1, 0, false)
	if err != nil {
		logger.Error(err.Error())
	}


	/********************************************************************/
	/* CONFIGURE MySql DB												*/
	/********************************************************************/
	// Load MySql configuration
	MySqlHost 				:= viper.GetString(`MySql.Host`)
	MySqlUserName 			:= viper.GetString(`MySql.UserName`)
	MySqlPassword			:= viper.GetString(`MySql.Password`)
	MySqlDatabase			:= viper.GetString(`MySql.Database`)
	MySqlMaxOpenConnections	:= viper.GetInt(`MySql.MaxOpenConnections`)
	MySqlMaxIdleConnections	:= viper.GetInt(`MySql.MaxIdleConnections`)

	//timeout := time.Duration(viper.GetInt("Context.Timeout")) * time.Second

	// Open a MySql infrastructure
	dbContext := repository.ConnectMySql(MySqlHost, MySqlUserName, MySqlPassword, MySqlDatabase, MySqlMaxOpenConnections, MySqlMaxIdleConnections)
	if dbContext == nil {
		os.Exit(1)
	}

	err = dbContext.Ping()
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Connected")
	}

	defer func() {
		err := dbContext.Close()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}()


	/********************************************************************/
	/* CONFIGURE ECHO													*/
	/********************************************************************/
	//e := echo.New()
	//e.Use(middleware.CORS())

	/********************************************************************/
	/* Redis												*/
	/********************************************************************/
	host := viper.GetString("Redis.Host")
	poolSize := viper.GetString("Redis.PoolSize")
	minIdleConns := viper.GetString("Redis.MinIdleConns")

	cacheManager := cache.CacheManager{}
	cacheManager.Init(host, util.ParseInt(poolSize), util.ParseInt(minIdleConns))
	pong, err := cacheManager.Client.Ping().Result()
	fmt.Println(pong, err)

	/********************************************************************/
	/* INITIALIZE MODULES												*/
	/********************************************************************/
	lottery.Initialize(rbChannel, dbContext, cacheManager, timeout)
	lottery.Execute()


}
