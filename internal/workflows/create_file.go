package workflows

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AB527/AWSEasyDeploy/internal/assets"
	"github.com/AB527/AWSEasyDeploy/internal/constants"
)

func RunSetupWorkflow() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	configPath := filepath.Join(wd, ".easy-deploy")

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config constants.DeployConfig
	if err := json.Unmarshal(configFile, &config); err != nil {
		return err
	}

	if updated, _ := json.MarshalIndent(config, "", "  "); updated != nil {
		_ = os.WriteFile(configPath, updated, 0644)
	}

	templatePath := fmt.Sprintf("workflows/%s.yaml", config.DeploymentSolution)
	raw, err := assets.Workflows.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("unable to load workflow template %s: %w", templatePath, err)
	}
	content := string(raw)

	inject := fmt.Sprintf(`on:
  push:
    branches:
      - %s

`, config.Branch)

	lines := strings.Split(content, "\n")
	var result []string

	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "name:") {
			result = append(result, inject)
			result = append(result, lines[i:]...)
			break
		}
		result = append(result, line)
	}

	replacements := map[string]string{
		"${{ secrets.AWS_ACCESS_KEY_ID }}":     config.AwsAccessKeyID,
		"${{ secrets.AWS_SECRET_ACCESS_KEY }}": config.AwsSecretAccessKey,
		"${{ secrets.AWS_REGION }}":            config.AwsRegion,
		"${{ secrets.AWS_S3_BUCKET }}":         config.AwsS3Bucket,
		"${{ secrets.AWS_EB_APP }}":            config.AwsEbApp,
		"${{ secrets.AWS_EB_ENV }}":            config.AwsEbEnv,
		"${{ secrets.BRANCH_NAME }}":           config.Branch,
	}

	for j := range result {
		for old, val := range replacements {
			if val == "" {
				continue
			}
			result[j] = strings.ReplaceAll(result[j], old, val)
		}
	}

	var outputDir, outputFile string
	switch config.DeploymentSolution {
	case "github":
		outputDir = filepath.Join(wd, ".github", "workflows")
		outputFile = "easy-deploy.yaml"
	case "gitlab":
		outputDir = wd
		outputFile = ".gitlab-ci.yml"
	default:
		return fmt.Errorf("unsupported deployment solution: %s", config.DeploymentSolution)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := os.WriteFile(outputPath, []byte(strings.Join(result, "\n")), 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Workflow written: %s\n", outputPath)
	return nil
}
