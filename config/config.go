package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWSConfigurator interface {
	LoadAWSConfig(profile string) (aws.Config, error)
}

type AWSConfigLoader struct {
	Cfg aws.Config
}

func (c *AWSConfigLoader) LoadAWSConfig(profile string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))

	if err != nil {
		return aws.Config{}, err
	}

	c.Cfg = cfg

	return cfg, nil
}
