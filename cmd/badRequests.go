/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/nobelsmith/go-fence/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// badRequestsCmd represents the badRequests command
var badRequestsCmd = &cobra.Command{
	Use:   "badRequests",
	Short: "Reads through log file and prints bad requests (4xx and 5xx)",
	Long: `Reads through a log file and prints requests that failed.
The output is grouped by ip address and sorted by the number of failed requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.BadRequests(viper.GetString("nginxlogfile"), viper.GetStringSlice("protectedips"), viper.GetStringSlice("forbiddenlocations"))
		if err != nil {
			log.Fatal(err)
		}
	},
}

// adds the badRequests command to the root command.
func init() {
	rootCmd.AddCommand(badRequestsCmd)
}
