package aws

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	elasticbeanstalktypes "github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateEBApplication(ctx context.Context, appName, description string) {
	client := getEBClient(ctx)

	_, err := client.CreateApplication(ctx, &elasticbeanstalk.CreateApplicationInput{
		ApplicationName: &appName,
		Description:     &description,
	})
	if err != nil {
		checkAndReportPermissionError(err)
		if strings.Contains(err.Error(), "Application already exists") {
			fmt.Printf("✅ EB application already exists: %s\n", appName)
			return
		}
		log.Fatalf("❌ Failed to create EB application: %v", err)
	}

	fmt.Printf("✅ EB application created: %s\n", appName)
}

func CreateEBEnvironment(ctx context.Context, appName, envName, solutionStack string) {
	client := getEBClient(ctx)

	if solutionStack == "" {
		log.Fatalf("❌ No solution stack resolved for environment %s", envName)
	}

	_, err := client.CreateEnvironment(ctx, &elasticbeanstalk.CreateEnvironmentInput{
		ApplicationName:   &appName,
		EnvironmentName:   &envName,
		SolutionStackName: &solutionStack,
		OptionSettings: []elasticbeanstalktypes.ConfigurationOptionSetting{
			{
				Namespace:  awsStringPtr("aws:autoscaling:launchconfiguration"),
				OptionName: awsStringPtr("IamInstanceProfile"),
				Value:      awsStringPtr(ebInstanceProfileName),
			},
		},
	})
	if err != nil {
		checkAndReportPermissionError(err)
		if strings.Contains(err.Error(), "Environment already exists") {
			fmt.Printf("✅ EB environment already exists: %s\n", envName)
			return
		}
		log.Fatalf("❌ Failed to create EB environment: %v", err)
	}

	fmt.Printf("✅ EB environment created: %s\n", envName)
}

func GetLatestSolutionStack(ctx context.Context, keyword string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := elasticbeanstalk.NewFromConfig(cfg)

	resp, err := client.ListAvailableSolutionStacks(ctx, &elasticbeanstalk.ListAvailableSolutionStacksInput{})
	if err != nil {
		return "", fmt.Errorf("failed to list solution stacks: %w", err)
	}

	for _, s := range resp.SolutionStacks {
		if strings.Contains(s, keyword) {
			return s, nil
		}
	}

	return "", fmt.Errorf("no solution stack found matching %q", keyword)
}

func UpdateEnvironmentVariables(ctx context.Context, envName string, vars map[string]string) error {
	client := getEBClient(ctx)

	var settings []elasticbeanstalktypes.ConfigurationOptionSetting
	for k, v := range vars {
		settings = append(settings, elasticbeanstalktypes.ConfigurationOptionSetting{
			Namespace:  awsStringPtr("aws:elasticbeanstalk:application:environment"),
			OptionName: awsStringPtr(k),
			Value:      awsStringPtr(v),
		})
	}

	_, err := client.UpdateEnvironment(ctx, &elasticbeanstalk.UpdateEnvironmentInput{
		EnvironmentName: awsStringPtr(envName),
		OptionSettings:  settings,
	})
	if err != nil {
		checkAndReportPermissionError(err)
	}
	return err
}

func CreateS3Bucket(ctx context.Context, bucketName string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to load AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		if strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") {
			fmt.Printf("✅ S3 bucket already exists: %s\n", bucketName)
			return nil
		}
		if strings.Contains(err.Error(), "BucketAlreadyExists") {
			return fmt.Errorf("bucket name not available: %s", bucketName)
		}
		return fmt.Errorf("failed to create S3 bucket: %w", err)
	}

	fmt.Printf("✅ S3 bucket created: %s\n", bucketName)
	return nil
}
