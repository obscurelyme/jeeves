package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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

	t.Run("should return an error attempting to log in", func(t *testing.T) {
		errorMessage := "Unable to login"
		err := Login(&LoginProvider{
			loginConfig: &MockConfig{
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

func TestLoginConfig(t *testing.T) {
	cfg := Config{}

	t.Run("should succeed when a valid aws credential and viper struct are passed in", func(t *testing.T) {
		expiresTime := time.Now()
		creds := aws.Credentials{
			AccessKeyID:     "123",
			SecretAccessKey: "123-secret",
			SessionToken:    "123-token",
			Expires:         expiresTime,
		}
		vCfg := viper.New()
		err := cfg.SyncSessionCredentials(creds, vCfg, &SyncSessionCredentialsInput{})

		if assert.NoError(t, err) {
			assert.Equal(t, creds.AccessKeyID, vCfg.GetString("default.aws_access_key_id"))
			assert.Equal(t, creds.SecretAccessKey, vCfg.GetString("default.aws_secret_access_key"))
			assert.Equal(t, creds.SessionToken, vCfg.GetString("default.aws_session_token"))
		}
	})

	t.Run("should succeed when a valid aws credential, viper struct, and custom input are passed in", func(t *testing.T) {
		expiresTime := time.Now()
		creds := aws.Credentials{
			AccessKeyID:     "123",
			SecretAccessKey: "123-secret",
			SessionToken:    "123-token",
			Expires:         expiresTime,
		}
		vCfg := viper.New()
		err := cfg.SyncSessionCredentials(creds, vCfg, &SyncSessionCredentialsInput{
			profile: "test",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, creds.AccessKeyID, vCfg.GetString("test.aws_access_key_id"))
			assert.Equal(t, creds.SecretAccessKey, vCfg.GetString("test.aws_secret_access_key"))
			assert.Equal(t, creds.SessionToken, vCfg.GetString("test.aws_session_token"))
		}
	})

	t.Run("should fail when the viper pointer is nil", func(t *testing.T) {
		expectedError := errors.New("vConfig cannot be nil")
		err := cfg.SyncSessionCredentials(aws.Credentials{}, nil, &SyncSessionCredentialsInput{})

		if assert.Error(t, err) {
			assert.Equal(t, expectedError, err)
		}
	})

}
