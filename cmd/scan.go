package cmd

import (
	"log"

	"github.com/nobelsmith/go-fence/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Reads through log file and bans users",
	Long: `Reads through a log file and will ban users if
	a they have tried to access a forbidden location.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.ReadConfig(CfgFile)
		if err != nil {
			log.Fatal(err)
		}
		err = pkg.Scan(viper.GetString("nginxlogfile"), viper.GetStringSlice("protectedips"), viper.GetStringSlice("forbiddenlocations"), DryRun)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// It adds the scan command to the root command.
func init() {
	rootCmd.AddCommand(scanCmd)
}
