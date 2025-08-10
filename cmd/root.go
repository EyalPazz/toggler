package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "toggler",
	Short: "A CLI for Toggl time tracking",
	Long:  `Toggler is a command line interface for the Toggl time tracking service.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("toggler")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("TOGGLER")
	viper.AutomaticEnv()
	
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}
}
