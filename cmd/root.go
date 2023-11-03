package cmd

import (
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile = ""
	version = false
	dryRun  = false
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
			log.Println("go-fence version 0.0.2")
		} else {
			log.Println(`Use "go-fence --help" for more information.`)
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (if left blank a default will be created at $HOME/.go-fence.yaml)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "output offending users without banning")

	//local flags
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "application version")
}

func initConfig() {
	if version {
		return
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
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

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no config file found")
			if cfgFile == "" {
				log.Println("using defaults")
				viper.SetDefault("NginxLogFile", "/var/log/nginx/access.log")
				viper.SetDefault("ProtectedIPs", []string{"192.168.*", "172.16.*", "10.*", "127*"})
				viper.SetDefault("ForbiddenLocations", []string{"wp-admin", "wp-login.php", ".aspx"})
				err := viper.SafeWriteConfig()
				if err != nil {
					log.Println("Error writing config file: ", err)
				}
			}
		} else {
			log.Println("problem with config file: ", err)
		}
	}
}
