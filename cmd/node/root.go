package node

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/3d0c/storage/pkg/config"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "node",
	Short: "Storage Node server",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// initConfig reads in config file
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		fmt.Println("Config file hasn't been provided")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else if cfgFile != "" {
		panic(fmt.Sprintf("Error reading config file from '%s' - %s", cfgFile, err))
	}

	if err := viper.Unmarshal(config.Node()); err != nil {
		panic(fmt.Sprintf("Failed to init config: %s", err))
	}

	// Override address if it's provided by ENV
	if address := os.Getenv("NODE_LISTEN"); address != "" {
		config.Node().Server.Address = address
	}

	// Override storage dir if it's provided by ENV
	if dir := os.Getenv("NODE_STORAGE"); dir != "" {
		config.Node().Saver.StorageDir = dir
	}
}
