package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/obscurelyme/jeeves/ini"
	"github.com/spf13/cobra"
)

var profile string
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to AWS",
	Long:  "Use Jeeves to login to AWS",
	Run:   loginToAws,
}

func init() {
	loginCmd.PersistentFlags().Bool("sso", true, "Login to AWS via SSO with IAM Identity Center (default \"true\")")
	loginCmd.PersistentFlags().Bool("assume-role", false, "Assume a role in AWS")
	loginCmd.PersistentFlags().Bool("session", false, "Generate a session token for AWS")
	loginCmd.PersistentFlags().StringVar(&profile, "profile", "default", "AWS Profile to login to")
	rootCmd.AddCommand(loginCmd)
}

func loginToAws(cmd *cobra.Command, args []string) {
	// Step 1, see if "jeeves" profile is set in .aws/config
	// Step 1.5, if no "jeeves", then make "jeeves"
	// Step 2, if "jeeves", then continue with SSO logon

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))

	if err != nil {
		// NOTE: try to make a profile
		cfg = awsConfigureSSO()
	}

	stsClient := sts.NewFromConfig(cfg)
	identity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	s := strings.Split(string(*identity.UserId), ":")

	stsClient.GetAccessKeyInfo(context.TODO(), &sts.GetAccessKeyInfoInput{
		AccessKeyId: &s[0],
	})
	stsClient.GetSessionToken(context.TODO(), &sts.GetSessionTokenInput{})

	if err != nil {
		log.Fatalln(err.Error())
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())

	if err != nil {
		log.Fatalln(err.Error())
	}

	if profile != "default" {
		AWSCredentials.Set(fmt.Sprintf("%s.aws_access_key_id", profile), creds.AccessKeyID)
		AWSCredentials.Set(fmt.Sprintf("%s.aws_secret_access_key", profile), creds.SecretAccessKey)
		AWSCredentials.Set(fmt.Sprintf("%s.aws_session_token", profile), creds.SessionToken)
	} else {
		AWSCredentials.Set("default.aws_access_key_id", creds.AccessKeyID)
		AWSCredentials.Set("default.aws_secret_access_key", creds.SecretAccessKey)
		AWSCredentials.Set("default.aws_session_token", creds.SessionToken)
	}

	var credsErr = ini.WriteIniFile(AWSCredentialsFilePath, AWSCredentials.AllSettings())

	if credsErr != nil {
		log.Fatalln(credsErr.Error())
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

	stsClient := sts.NewFromConfig(cfg)

	output, err := stsClient.AssumeRole(context.TODO(), &sts.AssumeRoleInput{
		RoleArn:         &RoleArn,
		RoleSessionName: &RoleSessionName,
		DurationSeconds: &DurationSeconds,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Print(*output.Credentials)
}

func awsConfigureSSO() aws.Config {
	awsCmd := exec.Command("aws", "configure", "sso", "--profile", profile)

	awsCmd.Stdin = os.Stdin
	awsCmd.Stdout = os.Stdout
	awsCmd.Stderr = os.Stderr

	err := awsCmd.Run()

	if err != nil {
		log.Fatalln(err.Error())
	}

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))

	return cfg
}
