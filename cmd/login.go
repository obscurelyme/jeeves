package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to AWS",
	Long:  "Use Jeeves to login to AWS",
	Run:   loginToAws,
}

func loginToAws(cmd *cobra.Command, args []string) {
	fmt.Println("You have logged into AWS!")
}
