package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	CfgFile = ""
	version = false
	DryRun  = false
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-fence",
	Short: "Simple nginx security service",
	Long: `go-fence is a simple nginx security service that
monitors nginx log files and ban users who access forbidden locations.
This application needs Nginx logs to be specifically formatted
as json and to be run on a system that is running iptables.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Println("go-fence version 0.0.3")
		} else {
			fmt.Println(`Use "go-fence --help" for more information.`)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// global flags
	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (if left blank a default will be created at $HOME/.go-fence.yaml)")
	rootCmd.PersistentFlags().BoolVar(&DryRun, "dry-run", false, "output offending users without banning")

	//local flags
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "application version")
}

func initConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		viper.SetConfigName(".go-fence") // name of config file (without extension)
		viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(home)        // path to look for the config file in
	}
}

// func initConfig() {
// 	if version {
// 		return
// 	}

// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := homedir.Dir()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		viper.SetConfigName(".go-fence") // name of config file (without extension)
// 		viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name
// 		viper.AddConfigPath(home)        // path to look for the config file in
// 	}

// 	if err := viper.ReadInConfig(); err != nil {
// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			// Config file not found; ignore error if desired
// 			fmt.Println("unable to find config file at ", cfgFile)
// 			if cfgFile == "" {
// 				fmt.Println(`no config file found, to create a default config file at $HOME/.go-fence.yaml by running "go-fence init"`)
// 			}
// 		} else {
// 			fmt.Println("problem with config file: ", err)
// 		}
// 	}
// }
