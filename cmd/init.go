package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a default config file at $HOME/.go-fence.yaml",
	Long: `Creates a default config file at $HOME/.go-fence.yaml with the following defaults:

NginxLogFile: /var/log/nginx/access.log
ProtectedIPs: ["192.168.*", "172.16.*", "10.*", "127*"]
ForbiddenLocations: ["wp-admin", "wp-login.php"]`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetDefault("NginxLogFile", "/var/log/nginx/access.log")
		viper.SetDefault("ProtectedIPs", []string{"192.168.*", "172.16.*", "10.*", "127*"})
		viper.SetDefault("ForbiddenLocations", []string{"wp-admin", "wp-login.php"})
		err := viper.SafeWriteConfig()
		if err != nil {
			log.Println("Error writing config file: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
