package ui

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"easy-deploy/internal/aws"
	"easy-deploy/internal/constants"
	"easy-deploy/internal/helpers"

	"github.com/manifoldco/promptui"
)

func RunAppDetails(ctx context.Context) {
	wd, _ := os.Getwd()
	configPath := filepath.Join(wd, ".easy-deploy")

	if _, err := os.Stat(configPath); err == nil {
		choicePrompt := promptui.Select{
			Label: "Configuration exists – reinitialize from it?",
			Items: []string{"Yes", "No"},
		}
		_, choice, err := choicePrompt.Run()
		if err != nil {
			return
		}

		if choice == "Yes" {
			data, err := os.ReadFile(configPath)
			if err != nil {
				fmt.Println("❌ Unable to read config")
				return
			}
			var cfg constants.DeployConfig
			if err := json.Unmarshal(data, &cfg); err != nil {
				fmt.Println("❌ Invalid config file")
				return
			}

			aws.ConnectToElasticBeanstalk(ctx)
			aws.CreateEBApplication(ctx, cfg.AwsEbApp, "")
			_ = aws.CreateS3Bucket(ctx, cfg.AwsS3Bucket)

			stack, err := resolveStack(ctx, cfg.Framework)
			if err != nil {
				fmt.Println("❌ Unable to resolve solution stack")
				return
			}

			if err := aws.EnsureInstanceProfile(ctx); err != nil {
				fmt.Println("❌ Failed to ensure instance profile")
				return
			}

			aws.CreateEBEnvironment(ctx, cfg.AwsEbApp, cfg.AwsEbEnv, stack)
			_ = aws.UpdateEnvironmentVariables(ctx, cfg.AwsEbEnv, map[string]string{
				"S3_BUCKET": cfg.AwsS3Bucket,
			})
			return
		}
	}

	projPrompt := promptui.Prompt{
		Label: "Project name",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("name cannot be empty")
			}
			return nil
		},
	}
	project, err := projPrompt.Run()
	if err != nil {
		return
	}

	frameworkPrompt := promptui.Select{
		Label: "Select your framework",
		Items: helpers.MapKeys(constants.FrameworkMap),
	}
	_, framework, err := frameworkPrompt.Run()
	if err != nil {
		return
	}

	deployPrompt := promptui.Select{
		Label: "Select code deployment solution",
		Items: helpers.MapKeys(constants.DeploymentMap),
	}
	_, deploymentSolution, err := deployPrompt.Run()
	if err != nil {
		return
	}

	branchPrompt := promptui.Prompt{
		Label:   "Enter branch name",
		Default: "main",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("branch name cannot be empty")
			}
			return nil
		},
	}
	branch, err := branchPrompt.Run()
	if err != nil {
		return
	}

	vpcPrompt := promptui.Select{
		Label: "Launch in VPC?",
		Items: []string{"Yes", "No"},
	}
	_, vpcChoice, err := vpcPrompt.Run()
	if err != nil {
		return
	}

	config := constants.DeployConfig{
		Framework:          NormalizeFramework(framework),
		DeploymentSolution: NormalizeDeploymentSolution(deploymentSolution),
		Branch:             branch,
		LaunchInVPC:        vpcChoice == "Yes",
		AwsRegion:          "us-east-1",
		AwsS3Bucket:        generateBucketName(project),
		AwsEbApp:           project,
		AwsEbEnv:           project + "-env",
	}

	confirmPrompt := promptui.Select{
		Label: "Are these details correct?",
		Items: []string{"Yes", "No"},
	}
	_, confirm, err := confirmPrompt.Run()
	if err != nil {
		return
	}

	if confirm != "Yes" {
		fmt.Println("Aborting. Please run `init` again to correct your inputs.")
		return
	}

	if err := SaveConfig(config); err != nil {
		fmt.Println("❌ Failed to save config")
		return
	}

	aws.ConnectToElasticBeanstalk(ctx)
	aws.CreateEBApplication(ctx, project, "")

	if err := aws.CreateS3Bucket(ctx, config.AwsS3Bucket); err != nil {
		fmt.Printf("❌ Could not create S3 bucket\n")
	}

	stack, err := resolveStack(ctx, config.Framework)
	if err != nil {
		fmt.Println("❌ Unable to resolve solution stack")
		return
	}

	if err := aws.EnsureInstanceProfile(ctx); err != nil {
		fmt.Println("❌ Failed to ensure instance profile")
		return
	}

	aws.CreateEBEnvironment(ctx, config.AwsEbApp, config.AwsEbEnv, stack)
	_ = aws.UpdateEnvironmentVariables(ctx, config.AwsEbEnv, map[string]string{
		"S3_BUCKET": config.AwsS3Bucket,
	})
}

func resolveStack(ctx context.Context, framework string) (string, error) {
	keyword, ok := constants.SolutionStackKeywordMap[framework]
	if !ok {
		return "", fmt.Errorf("no solution stack keyword defined for framework %q", framework)
	}

	stack, err := aws.GetLatestSolutionStack(ctx, keyword)
	if err != nil {
		if fallback, ok := constants.SolutionStackFallbackMap[framework]; ok {
			return fallback, nil
		}
		return "", fmt.Errorf("no fallback stack defined for framework %q", framework)
	}

	return stack, nil
}

func NormalizeFramework(input string) string {
	if val, ok := constants.FrameworkMap[input]; ok {
		return val
	}
	return strings.ToLower(input)
}

func NormalizeDeploymentSolution(input string) string {
	if val, ok := constants.DeploymentMap[input]; ok {
		return val
	}
	return strings.ToLower(input)
}

func SaveConfig(config constants.DeployConfig) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	filePath := filepath.Join(wd, ".easy-deploy")

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write json: %v", err)
	}

	fmt.Println("✅ Config saved")
	return nil
}

func generateBucketName(project string) string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%s-%s", project, hex.EncodeToString(b))
}
