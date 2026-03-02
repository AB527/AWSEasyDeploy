package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

const ebInstanceProfileName = "aws-elasticbeanstalk-ec2-role"

var ebInstanceProfilePolicies = []string{
	"arn:aws:iam::aws:policy/AWSElasticBeanstalkWebTier",
	"arn:aws:iam::aws:policy/AWSElasticBeanstalkWorkerTier",
	"arn:aws:iam::aws:policy/AWSElasticBeanstalkMulticontainerDocker",
}

func EnsureInstanceProfile(ctx context.Context) error {
	client := getIAMClient(ctx)

	if _, err := ensureRole(ctx, client); err != nil {
		return err
	}

	if _, err := ensureInstanceProfile(ctx, client); err != nil {
		return err
	}

	for i := 0; i < 20; i++ {
		_, err := client.GetInstanceProfile(ctx, &iam.GetInstanceProfileInput{
			InstanceProfileName: awsStringPtr(ebInstanceProfileName),
		})
		if err == nil {
			fmt.Println("✅ Instance profile ready")
			time.Sleep(10 * time.Second)
			return nil
		}
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("instance profile %s did not become available after 100 seconds", ebInstanceProfileName)
}

func ensureRole(ctx context.Context, client *iam.Client) (bool, error) {
	_, err := client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: awsStringPtr(ebInstanceProfileName),
	})
	if err == nil {
		fmt.Println("✅ IAM role already exists")
		return false, nil
	}
	if !strings.Contains(err.Error(), "NoSuchEntity") {
		return false, fmt.Errorf("failed to check IAM role: %w", err)
	}

	trustPolicy := map[string]any{
		"Version": "2012-10-17",
		"Statement": []map[string]any{
			{
				"Effect":    "Allow",
				"Principal": map[string]any{"Service": "ec2.amazonaws.com"},
				"Action":    "sts:AssumeRole",
			},
		},
	}
	trustJSON, err := json.Marshal(trustPolicy)
	if err != nil {
		return false, fmt.Errorf("failed to marshal trust policy: %w", err)
	}

	_, err = client.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 awsStringPtr(ebInstanceProfileName),
		AssumeRolePolicyDocument: awsStringPtr(string(trustJSON)),
	})
	if err != nil {
		return false, fmt.Errorf("failed to create IAM role: %w", err)
	}
	fmt.Println("✅ IAM role created")

	for _, policyArn := range ebInstanceProfilePolicies {
		_, err := client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
			RoleName:  awsStringPtr(ebInstanceProfileName),
			PolicyArn: awsStringPtr(policyArn),
		})
		if err != nil {
			return false, fmt.Errorf("failed to attach policy %s: %w", policyArn, err)
		}
	}
	fmt.Println("✅ IAM policies attached")

	return true, nil
}

func ensureInstanceProfile(ctx context.Context, client *iam.Client) (bool, error) {
	_, err := client.GetInstanceProfile(ctx, &iam.GetInstanceProfileInput{
		InstanceProfileName: awsStringPtr(ebInstanceProfileName),
	})
	if err == nil {
		fmt.Println("✅ Instance profile already exists")
		return false, nil
	}
	if !strings.Contains(err.Error(), "NoSuchEntity") {
		return false, fmt.Errorf("failed to check instance profile: %w", err)
	}

	_, err = client.CreateInstanceProfile(ctx, &iam.CreateInstanceProfileInput{
		InstanceProfileName: awsStringPtr(ebInstanceProfileName),
	})
	if err != nil {
		return false, fmt.Errorf("failed to create instance profile: %w", err)
	}

	_, err = client.AddRoleToInstanceProfile(ctx, &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: awsStringPtr(ebInstanceProfileName),
		RoleName:            awsStringPtr(ebInstanceProfileName),
	})
	if err != nil {
		return false, fmt.Errorf("failed to attach role to instance profile: %w", err)
	}

	fmt.Println("✅ Instance profile created")
	return true, nil
}
