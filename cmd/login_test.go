package cmd

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type MockConfig struct {
	Configurator
	output                     aws.Config
	err                        error
	ssoLoginErr                error
	configureSSOErr            error
	ssoSessionCredentialsError error
	ssoCredenials              aws.Credentials
}

func (c *MockConfig) LoadAWSConfig(profile string) (aws.Config, error) {
	return c.output, c.err
}

func (c *MockConfig) SSOLogin() error {
	return c.ssoLoginErr
}

func (c *MockConfig) ConfigureSSO() error {
	return c.configureSSOErr
}

func (c *MockConfig) GetSSOSessionCredentials(cfg aws.Config) (aws.Credentials, error) {
	return c.ssoCredenials, c.ssoSessionCredentialsError
}

func TestLogin(t *testing.T) {
	t.Run("should be successful", func(t *testing.T) {
		err := Login(&LoginProvider{
			loginConfig: &MockConfig{
				output: aws.Config{},
				err:    nil,
				ssoCredenials: aws.Credentials{
					CanExpire: true,
					Expires:   time.Now().AddDate(1, 0, 0), // NOTE: a year into the future to guarentee that the creds are not expired
				},
			},
			profile: "default",
		})

		if err != nil {
			t.Error("Login to AWS returned an error")
		}
	})
}
