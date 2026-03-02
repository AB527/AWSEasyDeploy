package aws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
)

var requiredPermissions = []string{
	"elasticbeanstalk:DescribeApplications",
	"elasticbeanstalk:DescribeEnvironments",
	"elasticbeanstalk:DescribeApplicationVersions",
	"elasticbeanstalk:DescribeConfigurationSettings",
}

func ConnectToElasticBeanstalk(ctx context.Context) {
	client := getEBClient(ctx)
	_, err := client.DescribeApplications(ctx, &elasticbeanstalk.DescribeApplicationsInput{})
	if err != nil {
		checkAndReportPermissionError(err)
		log.Fatalf("❌ Failed to connect to Elastic Beanstalk: %v", err)
	}
	fmt.Println("✅ Connected to Elastic Beanstalk")
}

func getIAMClient(ctx context.Context) *iam.Client {
	cfg, err := loadAWSConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}
	return iam.NewFromConfig(cfg)
}

func checkIAMPermissions(ctx context.Context, actions []string) {
	iamClient := getIAMClient(ctx)

	result, err := iamClient.SimulatePrincipalPolicy(ctx, &iam.SimulatePrincipalPolicyInput{
		PolicySourceArn: nil,
		ActionNames:     actions,
	})
	if err != nil {
		log.Fatalf("❌ Failed to simulate IAM policy: %v", err)
	}

	for _, res := range result.EvaluationResults {
		if res.EvalDecision == iamtypes.PolicyEvaluationDecisionTypeAllowed {
			fmt.Printf("✅ %s\n", *res.EvalActionName)
		} else {
			fmt.Printf("❌ %s\n", *res.EvalActionName)
		}
	}
}
