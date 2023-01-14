package main

import (
	"butterfly/butterfly"
	"butterfly/config"
	"github.com/spf13/viper"
	"log"
)

func main() {
	// check the db configuration (panic on failure)
	dbConfig := viper.New()
	dbConfig.AddConfigPath("./")
	dbConfig.SetConfigName("config")
	dbConfig.SetConfigType("yml")

	if err := dbConfig.ReadInConfig(); err != nil {
		log.Fatalf("check configuration failed: %v", err)
	}
	var pgConfig config.PgConfig

	if err := dbConfig.Unmarshal(&pgConfig); err != nil {
		log.Fatalf("get DB configuration failed: %v", err)
	}
	pgConfig.Connect()

	butterfly.Constructor(&pgConfig)
	generater := butterfly.NewWithDbConfig()
	generater.Commit()

	//TODO insert or update db table for butterfly
	//TODO create the instance of butterfly to generate id

}
