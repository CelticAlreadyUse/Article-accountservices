package cmd

import (
	"log"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gb-2-api-account-service",
	Short: "api service for account api gateway",
}
func initConfig(){
	config.InitConfig()
}
func Execute(){
	initConfig()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}