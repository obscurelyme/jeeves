package faas

import (
	"errors"

	"github.com/spf13/cobra"
)

var updateFaasCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the function code of a FaaS resource",
	Long: `Updates the function code of a FaaS resource. 
for Java resources you would upload the shaded jar, for everything else
please upload a zip file`,
	RunE: updateFaasCmdHandler,
}

func updateFaasCmdHandler(cmd *cobra.Command, args []string) error {
	return errors.New("this function is not implemented")
}
