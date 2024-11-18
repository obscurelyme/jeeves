package cmd

import (
	"log"

	"github.com/icza/gox/osx"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open your browser to AWS",
	Long:  "Use Jeeves to open your default browser to AWS SSO start url",
	Run:   openAWSStart,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func openAWSStart(cmd *cobra.Command, args []string) {
	startUrl := JeevesConfig.GetString("AWS.SSO.Start")

	if startUrl == "" {
		log.Fatalln("No Start URL present in Jeeves config file!")
	}

	err := osx.OpenDefault(startUrl)

	if err != nil {
		log.Fatalln(err)
	}
}
