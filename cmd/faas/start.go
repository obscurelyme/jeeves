package faas

import (
	"errors"

	"github.com/spf13/cobra"
)

var startFaasCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a local FaaS resource",
	Long:  "Starts a FaaS resource locally, using docker",
	RunE:  startFaasCmdHandler,
}

func startFaasCmdHandler(cmd *cobra.Command, args []string) error {
	return errors.New("this function is not implemented")
}
