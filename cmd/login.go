package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/obscurelyme/jeeves/config"
	"github.com/obscurelyme/jeeves/ini"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	SSOLogin() error
	ConfigureSSO() error
	GetSSOSessionCredentials(cfg aws.Config) (aws.Credentials, error)
	WriteSessionCredentials(filename string, vConfig *viper.Viper) error
	SyncSessionCredentials(creds aws.Credentials, vConfig *viper.Viper, options *SyncSessionCredentialsInput) error
}

type Config struct{}

func (c *Config) SSOLogin() error {
	awsCmd := exec.Command("aws", "sso", "login", "--profile", profile)

	awsCmd.Stdin = os.Stdin
	awsCmd.Stdout = os.Stdout
	awsCmd.Stderr = os.Stderr

	return awsCmd.Run()
}

func (c *Config) GetSSOSessionCredentials(cfg aws.Config) (aws.Credentials, error) {
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

func (c *Config) ConfigureSSO() error {
	awsCmd := exec.Command("aws", "configure", "sso", "--profile", profile)

	awsCmd.Stdin = os.Stdin
	awsCmd.Stdout = os.Stdout
	awsCmd.Stderr = os.Stderr

	return awsCmd.Run()
}

type SyncSessionCredentialsInput struct {
	profile string
}

func (c *Config) SyncSessionCredentials(creds aws.Credentials, vConfig *viper.Viper, options *SyncSessionCredentialsInput) error {
	if vConfig == nil {
		return errors.New("vConfig cannot be nil")
	}

	if options.profile != "default" && options.profile != "" {
		vConfig.Set(fmt.Sprintf("%s.aws_access_key_id", options.profile), creds.AccessKeyID)
		vConfig.Set(fmt.Sprintf("%s.aws_secret_access_key", options.profile), creds.SecretAccessKey)
		vConfig.Set(fmt.Sprintf("%s.aws_session_token", options.profile), creds.SessionToken)
		vConfig.Set(fmt.Sprintf("%s.aws_expires", options.profile), creds.Expires.String())
	} else {
		vConfig.Set("default.aws_access_key_id", creds.AccessKeyID)
		vConfig.Set("default.aws_secret_access_key", creds.SecretAccessKey)
		vConfig.Set("default.aws_session_token", creds.SessionToken)
		vConfig.Set("default.aws_expires", creds.Expires.String())
	}

	return nil
}

func (c *Config) WriteSessionCredentials(filename string, vConfig *viper.Viper) error {
	return ini.WriteIniFile(filename, vConfig.AllSettings())
}

type LoginProvider struct {
	loginConfig Configurator
	profile     string
}

func Login(provider *LoginProvider) error {
	// Check if we have a valid ~/.aws/config + profile
	loader := config.AWSConfigLoader{}
	cfg, err := loader.LoadAWSConfig(provider.profile)

	if err != nil {
		nerr := provider.loginConfig.ConfigureSSO()
		if nerr != nil {
			return nerr
		}
		// NOTE: SSO configured, at this point the user is logged in
		fmt.Println("AWS Login Successful!")
		return nil
	}

	// Check if we are still logged in, if not log in and refetch the credentials
	creds, err := provider.loginConfig.GetSSOSessionCredentials(cfg)

	if err != nil {
		return err
	}

	if creds.Expired() {
		// Login again
		nerr := provider.loginConfig.SSOLogin()
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
		loginConfig: &Config{},
		profile:     profile,
	}

	return Login(&provider)
}
