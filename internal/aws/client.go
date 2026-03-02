package aws

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
)

func loadAWSConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}

func getEBClient(ctx context.Context) *elasticbeanstalk.Client {
	cfg, err := loadAWSConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}
	return elasticbeanstalk.NewFromConfig(cfg)
}

func awsString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func awsStringPtr(s string) *string {
	return &s
}

func checkAndReportPermissionError(err error) {
	if err == nil {
		return
	}
	lower := strings.ToLower(err.Error())
	if strings.Contains(lower, "access denied") || strings.Contains(lower, "403") || strings.Contains(lower, "not authorized") {
		fmt.Println("❌ AWS access denied. Ensure your IAM credentials have the required Elastic Beanstalk permissions:")
		fmt.Println("   https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/iam-roles.html#iam-user-permissions")
	}
}
