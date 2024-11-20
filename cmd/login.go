package cmd

import (
	"context"
	"fmt"
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
	RunE:  loginToAws,
}

func init() {
	loginCmd.PersistentFlags().Bool("sso", true, "Login to AWS via SSO with IAM Identity Center (default \"true\")")
	loginCmd.PersistentFlags().Bool("assume-role", false, "Assume a role in AWS")
	loginCmd.PersistentFlags().Bool("session", false, "Generate a session token for AWS")
	loginCmd.PersistentFlags().StringVar(&profile, "profile", "default", "AWS Profile to login to")
	rootCmd.AddCommand(loginCmd)
}

type Configurator interface {
	LoadAWSConfig(profile string) (aws.Config, error)
	SSOLogin() error
	ConfigureSSO() error
	GetSSOSessionCredentials() (aws.Credentials, error)
	WriteSessionCredentials(creds aws.Credentials) error
}

type Config struct {
	awsConfig aws.Config
}

func (c *Config) LoadAWSConfig(profile string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))

	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

func (c *Config) SSOLogin() error {
	awsCmd := exec.Command("aws", "sso", "login", "--profile", profile)

	awsCmd.Stdin = os.Stdin
	awsCmd.Stdout = os.Stdout
	awsCmd.Stderr = os.Stderr

	return awsCmd.Run()
}

func (c *Config) GetSSOSessionCredentials() (aws.Credentials, error) {
	ssoClient := sso.NewFromConfig(c.awsConfig)
	ssoOidcClient := ssooidc.NewFromConfig(c.awsConfig)

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

func (c *Config) ConfigureSSO() error {
	awsCmd := exec.Command("aws", "configure", "sso", "--profile", profile)

	awsCmd.Stdin = os.Stdin
	awsCmd.Stdout = os.Stdout
	awsCmd.Stderr = os.Stderr

	return awsCmd.Run()
}

func (c *Config) WriteSessionCredentials(creds aws.Credentials) error {
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

	return ini.WriteIniFile(AWSCredentialsFilePath, AWSCredentials.AllSettings())
}

type LoginProvider struct {
	config  Configurator
	profile string
}

func Login(provider *LoginProvider) error {
	// Check if we have a valid ~/.aws/config + profile
	_, err := provider.config.LoadAWSConfig(provider.profile)

	if err != nil {
		nerr := provider.config.ConfigureSSO()
		if nerr != nil {
			return nerr
		}
		// NOTE: SSO configured, start over
		// return Login(provider)
	}

	// Check if we are still logged in, if not log in and refetch the credentials
	creds, err := provider.config.GetSSOSessionCredentials()

	if err != nil {
		return err
	}

	if creds.Expired() {
		// Login again
		nerr := provider.config.SSOLogin()
		if nerr != nil {
			return nerr
		}
		// NOTE: We are now logged in, start over
		// return Login(provider)
	}

	// NOTE: This is a VERY optional step in the SSO process btw...
	// only relevant if we want to copy the credentials over to
	// a docker container or something.
	// Write the creds to the ~/.aws/credentials shared file
	// err = provider.config.WriteSessionCredentials(creds)

	// if err != nil {
	// 	return err
	// }

	fmt.Println("AWS Login Successful!")
	return nil
}

func loginToAws(cmd *cobra.Command, args []string) error {
	var provider = LoginProvider{
		config:  &Config{},
		profile: profile,
	}

	return Login(&provider)
}
