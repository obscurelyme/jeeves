package cmd

import (
	"errors"
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

func (c *MockConfig) GetSSOSessionCredentials() (aws.Credentials, error) {
	return c.ssoCredenials, c.ssoSessionCredentialsError
}

func TestLogin(t *testing.T) {
	t.Run("should be successful", func(t *testing.T) {
		err := Login(&LoginProvider{
			config: &MockConfig{
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

	t.Run("should return an error attempting to read the .aws/config file", func(t *testing.T) {
		errorMessage := "Could not read or write .aws/config file"
		err := Login(&LoginProvider{
			config: &MockConfig{
				output:          aws.Config{},
				err:             errors.New("Test Error"),
				configureSSOErr: errors.New(errorMessage),
				ssoCredenials: aws.Credentials{
					CanExpire: true,
					Expires:   time.Now().AddDate(1, 0, 0), // NOTE: a year into the future to guarentee that the creds are not expired
				},
			},
			profile: "default",
		})

		if err.Error() != errorMessage {
			t.Errorf("Expected: %s, Actual: %s", errorMessage, err.Error())
		}
	})

	t.Run("should return an error attempting to read the .aws/sso/cache", func(t *testing.T) {
		errorMessage := "Unable to read .aws/sso/cache"
		err := Login(&LoginProvider{
			config: &MockConfig{
				output:                     aws.Config{},
				ssoSessionCredentialsError: errors.New(errorMessage),
				ssoCredenials: aws.Credentials{
					CanExpire: true,
					Expires:   time.Now().AddDate(1, 0, 0), // NOTE: a year into the future to guarentee that the creds are not expired
				},
			},
			profile: "default",
		})

		if err.Error() != errorMessage {
			t.Errorf("Expected: %s, Actual: %s", errorMessage, err.Error())
		}
	})

	t.Run("should return an error attempting to log in", func(t *testing.T) {
		errorMessage := "Unable to login"
		err := Login(&LoginProvider{
			config: &MockConfig{
				output:      aws.Config{},
				ssoLoginErr: errors.New(errorMessage),
				ssoCredenials: aws.Credentials{
					CanExpire: true,
					Expires:   time.Now(), // NOTE: a year into the future to guarentee that the creds are not expired
				},
			},
			profile: "default",
		})

		if err.Error() != errorMessage {
			t.Errorf("Expected: %s, Actual: %s", errorMessage, err.Error())
		}
	})
}
