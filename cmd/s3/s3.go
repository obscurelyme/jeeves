package s3

import "github.com/spf13/cobra"

var S3RootCmd = &cobra.Command{
	Use:   "s3",
	Short: "Execute S3 commands",
	Long:  "Configure, Create, List, and Delete S3 Buckets within AWS",
}

func init() {
	// TODO: Create these
	// s3Cmd.AddCommand(listCmd)
	// s3Cmd.AddCommand(createCmd)
	// s3Cmd.AddCommand(deleteCmd)
	// s3Cmd.AddCommand(updateCmd)
}
