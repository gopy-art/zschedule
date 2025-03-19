package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ConfigFile string
	Logger     string
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
	rootCmd.Flags().StringVarP(&Logger, "logger", "l", "stdout", "set app logger type , stdout or file")
	rootCmd.Flags().BoolVarP(&Version, "version", "v", false, "zschedule version")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
