package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to AWS",
	Long:  "Use Jeeves to login to AWS",
	Run:   loginToAws,
}

func init() {
	loginCmd.PersistentFlags().Bool("assume-role", false, "Assume a role in AWS")
	loginCmd.PersistentFlags().Bool("session", true, "Generate a session token for AWS (default \"true\")")
	rootCmd.AddCommand(loginCmd)
}

func loginToAws(cmd *cobra.Command, args []string) {
	shouldAssumeRole, _ := cmd.Flags().GetBool("assume-role")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedCredentialsFiles([]string{
		"~/.aws/credentials",
		"jeeves",
	}))
	if err != nil {
		// NOTE: Something went wrong reading .aws/config
		log.Fatal(err)
	} else {

		if shouldAssumeRole {
			assumeRole(cfg)
			return
		}

		getSessionToken(cfg)
	}

	log.Print("AWS Login Successful!")
}

// Generate a new session token for use of AWS resources
// using the current logged in account.
func getSessionToken(cfg aws.Config) {
	var DurationSeconds int32 = 3600

	output, err := sts.NewFromConfig(cfg).GetSessionToken(context.TODO(), &sts.GetSessionTokenInput{
		DurationSeconds: &DurationSeconds,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Print(*output.Credentials)

	home, _ := os.UserHomeDir()
	file, err := os.OpenFile(fmt.Sprintf("%v/.aws/credentials", home), os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	defer file.Close()
	str := []string{
		fmt.Sprintf("aws_access_key_id: %v\n", *output.Credentials.AccessKeyId),
		fmt.Sprintf("aws_secret_access_key: %v\n", *output.Credentials.SecretAccessKey),
		fmt.Sprintf("session_token: %v\n", *output.Credentials.SessionToken),
	}
	log.Print(str)
	// data := []byte(strings.Join(str, ""))

	// _, err = file.Write(data)
	// if err != nil {
	// 	panic(err)
	// }
}

// Assumes role, typically from another AWS account
func assumeRole(cfg aws.Config) {
	var DurationSeconds int32 = 3600
	var RoleArn string = "arn"
	RoleSessionName := fmt.Sprintf("JeevesSessionAssumedRole%v", "TempRoleName")

	output, err := sts.NewFromConfig(cfg).AssumeRole(context.TODO(), &sts.AssumeRoleInput{
		RoleArn:         &RoleArn,
		RoleSessionName: &RoleSessionName,
		DurationSeconds: &DurationSeconds,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Print(*output.Credentials)
}
