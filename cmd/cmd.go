package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ConfigFile string
	EnvFile    string
	Logger     string
	Type       string
	Version    bool
)

var rootCmd = &cobra.Command{
	Use:   "zschedule",
	Short: "simple and fast scheduler over http and commad line written in golang",
	Run: func(cmd *cobra.Command, args []string) {
		if Version {
			fmt.Println(AppVersion)
			os.Exit(0)
		}
	},
}

var AppVersion string = "v0.1.0"

// initial cli commands
func init() {
	rootCmd.Flags().StringVarP(&ConfigFile, "config", "c", "", "set the path of the config file. (for example : /var/config/json)")
	rootCmd.Flags().StringVarP(&Type, "type", "t", "", "set the type of the module. (api, cli)")
	rootCmd.Flags().StringVar(&EnvFile, "env", "", "set the .env file path for api server configuration.")
	rootCmd.Flags().StringVarP(&Logger, "logger", "l", "stdout", "set app logger type , stdout or file")
	rootCmd.Flags().BoolVarP(&Version, "version", "v", false, "zschedule version")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Validate() {
	if Type != "cli" && Type != "api" {
		fmt.Println("invalid type\nthe module type should be api or cli")
		os.Exit(1)
	}
}