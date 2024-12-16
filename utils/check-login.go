package utils

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/obscurelyme/jeeves/config"
)

var ErrNotLoggedIn error = errors.New("you need to login into AWS first, please run \"jeeves login\" then retry")

func CheckAWSLogin() (bool, error) {
	loader := config.AWSConfigLoader{}
	cfg, err := loader.LoadAWSConfig("default")
	if err != nil {
		return false, err
	}

	stsClient := sts.NewFromConfig(cfg)
	_, err = stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return false, err
	}

	return true, nil
}
