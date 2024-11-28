package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	jeevesConfig "github.com/obscurelyme/jeeves/config"
)

type AWSSSOConfig struct {
	SSOSessionName string
	SSOAccountId   string
	SSORoleName    string
	SSOStartURL    string
}

func awsSSOConfig(profile string) AWSSSOConfig {
	ssoSessionName := AWSConfig.GetString(fmt.Sprintf("%s.sso_session", profile))

	return AWSSSOConfig{
		SSOSessionName: ssoSessionName,
		SSOAccountId:   AWSConfig.GetString(fmt.Sprintf("%s.sso_account_id", profile)),
		SSORoleName:    AWSConfig.GetString(fmt.Sprintf("%s.sso_role_name", profile)),
		SSOStartURL:    AWSConfig.GetString(fmt.Sprintf("sso-session %s.sso_start_url", ssoSessionName)),
	}
}

// Returns the current session credentials from a valid SSO session
// **The client MUST be logged in via `jeeves login` before this function**
// **will return valid credentials**
func GetSSOSessionCredentials(profile string) (aws.Credentials, error) {
	loader := jeevesConfig.AWSConfigLoader{}
	cfg, err := loader.LoadAWSConfig(profile)

	if err != nil {
		return aws.Credentials{}, err
	}

	ssoClient := sso.NewFromConfig(cfg)
	ssoOidcClient := ssooidc.NewFromConfig(cfg)
	ssoAwsConfig := awsSSOConfig(profile)

	tokenPath, err := ssocreds.StandardCachedTokenFilepath(ssoAwsConfig.SSOSessionName)
	if err != nil {
		return aws.Credentials{}, err
	}

	var provider aws.CredentialsProvider
	provider = ssocreds.New(
		ssoClient,
		ssoAwsConfig.SSOAccountId,
		ssoAwsConfig.SSORoleName,
		ssoAwsConfig.SSOSessionName,
		func(options *ssocreds.Options) {
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
