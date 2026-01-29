package workflows

import (
	"easy-deploy-cli/internal/constants"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	assetPath := filepath.Join(wd, "assets", strings.ToLower(config.DeploymentSolution)+".yaml")
	content, err := os.ReadFile(assetPath)
	if err != nil {
		return fmt.Errorf("unsupported deployment solution")
	}

	inject := fmt.Sprintf(`on:
  push:
    branches:
      - %s

`, config.Branch)

	lines := strings.Split(string(content), "\n")
	var result []string

	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "name:") {
			result = append(result, inject)
			result = append(result, lines[i:]...)
			break
		}
		result = append(result, line)
	}

	var outputDir, outputFile string
	switch strings.ToLower(config.DeploymentSolution) {
	case "github":
		outputFile = "easy-deploy.yaml"
		outputDir = filepath.Join(wd, ".github", "workflows")
	case "gitlab":
		outputFile = "gitlab-ci.yaml"
		outputDir = filepath.Join(wd)
	default:
		outputDir = wd
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, outputFile)
	return os.WriteFile(outputPath, []byte(strings.Join(result, "\n")), 0644)
}
