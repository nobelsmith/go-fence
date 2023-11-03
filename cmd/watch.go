package cmd

import (
	"log"

	"github.com/nobelsmith/go-fence/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Mimicks tail -f funcionality on log file and bans users",
	Long: `Follows a log file and will survive things like log rotate.
If a user hits a forbidden location they will be banned.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.Watch(viper.GetString("nginxlogfile"), viper.GetStringSlice("protectedips"), viper.GetStringSlice("forbiddenlocations"), dryRun)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// It adds the watch command to the root command.
func init() {
	rootCmd.AddCommand(watchCmd)
}
