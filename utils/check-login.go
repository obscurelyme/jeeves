package utils

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/obscurelyme/jeeves/config"
)

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

	return false, errors.New("Hi")

	return true, nil
}
