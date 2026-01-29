package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"easy-deploy-cli/internal/aws"
	"easy-deploy-cli/internal/constants"
	"easy-deploy-cli/internal/helpers"

	"github.com/manifoldco/promptui"
)

// RunAppDetails handles the full interactive flow
func RunAppDetails(ctx context.Context) {
	// 1. Framework
	frameworkPrompt := promptui.Select{
		Label: "Select your framework",
		Items: helpers.MapKeys(constants.FrameworkMap),
	}
	_, framework, err := frameworkPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// 2. Deployment solution
	deployPrompt := promptui.Select{
		Label: "Select code deployment solution",
		Items: helpers.MapKeys(constants.DeploymentMap),
	}
	_, deploymentSolution, err := deployPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// 3. Branch
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
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// 4. VPC
	vpcPrompt := promptui.Select{
		Label: "Launch in VPC?",
		Items: []string{"Yes", "No"},
	}
	_, vpcChoice, err := vpcPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	launchInVPC := vpcChoice == "Yes"

	// Create config with normalized values
	config := constants.DeployConfig{
		Framework:          NormalizeFramework(framework),
		DeploymentSolution: NormalizeDeploymentSolution(deploymentSolution),
		Branch:             branch,
		LaunchInVPC:        launchInVPC,
	}

	// 5. Confirm details
	fmt.Println("\nPlease confirm your details:")
	fmt.Printf("Framework: %s\n", config.Framework)
	fmt.Printf("Deployment Solution: %s\n", config.DeploymentSolution)
	fmt.Printf("Branch Name: %s\n", config.Branch)
	fmt.Printf("Launch in VPC: %t\n", config.LaunchInVPC)

	confirmPrompt := promptui.Select{
		Label: "Are these details correct?",
		Items: []string{"Yes", "No"},
	}
	_, confirm, err := confirmPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if confirm != "Yes" {
		fmt.Println("Aborting. Please run `init` again to correct your inputs.")
		return
	}

	// 6. Save config
	if err := SaveConfig(config); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}

	// 7. AWS Deployment logic
	aws.ConnectToElasticBeanstalk()
	aws.CreateEBApplication(ctx, "MyApp", "My first Elastic Beanstalk application")
}

// NormalizeFramework converts human-readable to normalized
func NormalizeFramework(input string) string {
	if val, ok := constants.FrameworkMap[input]; ok {
		return val
	}
	return strings.ToLower(input)
}

// NormalizeDeploymentSolution converts human-readable to normalized
func NormalizeDeploymentSolution(input string) string {
	if val, ok := constants.DeploymentMap[input]; ok {
		return val
	}
	return strings.ToLower(input)
}

// SaveConfig writes JSON file in ui/app-details
func SaveConfig(config constants.DeployConfig) error {
	// directory where CLI is executed
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

	fmt.Println("✅ Config saved at:", filePath)
	return nil
}
