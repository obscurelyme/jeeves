package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

const VERSION = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Jeeves",
	Long:  "All software has versions. This is Jeeves's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Jeeves v%v\n", VERSION)
	},
}
