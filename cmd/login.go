package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
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
	var creds aws.Credentials

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		// NOTE: try to make a profile, and log in
		awsConfigureSSO()
	} else {
		// confirm we are logged in
		var e error
		creds, e = confirmSSOLogin(cfg)
		if e != nil || creds.Expired() {
			// log in
			awsConfigureSSO()
			creds, _ = confirmSSOLogin(cfg)
		}
	}

	if profile != "default" {
		AWSCredentials.Set(fmt.Sprintf("%s.aws_access_key_id", profile), creds.AccessKeyID)
		AWSCredentials.Set(fmt.Sprintf("%s.aws_secret_access_key", profile), creds.SecretAccessKey)
		AWSCredentials.Set(fmt.Sprintf("%s.aws_session_token", profile), creds.SessionToken)
		AWSCredentials.Set(fmt.Sprintf("%s.aws_expires", profile), creds.Expires.String())
	} else {
		AWSCredentials.Set("default.aws_access_key_id", creds.AccessKeyID)
		AWSCredentials.Set("default.aws_secret_access_key", creds.SecretAccessKey)
		AWSCredentials.Set("default.aws_session_token", creds.SessionToken)
		AWSCredentials.Set("default.aws_expires", creds.Expires.String())
	}

	var credsErr = ini.WriteIniFile(AWSCredentialsFilePath, AWSCredentials.AllSettings())

	if credsErr != nil {
		log.Fatalln(credsErr.Error())
	}

	log.Print("AWS Login Successful!")
}

// Confirms the current logon state for SSO, if the user is logged in,
// their credentials should exist, be valid and not be expired.
func confirmSSOLogin(cfg aws.Config) (aws.Credentials, error) {
	ssoClient := sso.NewFromConfig(cfg)
	ssoOidcClient := ssooidc.NewFromConfig(cfg)

	ssoSessionName := AWSConfig.GetString(fmt.Sprintf("%s.sso_session", profile))
	ssoAccountId := AWSConfig.GetString(fmt.Sprintf("%s.sso_account_id", profile))
	ssoRoleName := AWSConfig.GetString(fmt.Sprintf("%s.sso_role_name", profile))
	ssoStartUrl := AWSConfig.GetString(fmt.Sprintf("sso-session %s.sso_start_url", ssoSessionName))

	tokenPath, err := ssocreds.StandardCachedTokenFilepath(ssoSessionName)
	if err != nil {
		return aws.Credentials{}, err
	}

	var provider aws.CredentialsProvider
	provider = ssocreds.New(ssoClient, ssoAccountId, ssoRoleName, ssoStartUrl, func(options *ssocreds.Options) {
		options.SSOTokenProvider = ssocreds.NewSSOTokenProvider(ssoOidcClient, tokenPath)
	})

	// Wrap the provider with aws.CredentialsCache to cache the credentials until their expire time
	provider = aws.NewCredentialsCache(provider)

	credentials, err := provider.Retrieve(context.TODO())
	if err != nil {
		return aws.Credentials{}, err
	}

	return credentials, nil
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
