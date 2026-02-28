package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"easy-deploy-cli/internal/aws"
	"easy-deploy-cli/internal/constants"

	"github.com/spf13/cobra"
)

var pushEnvCmd = &cobra.Command{
	Use:   "push-env [file]",
	Short: "Push environment variables to configured EB environment",
	Long: `Read a file of KEY=VALUE pairs and apply them to the EB environment
named in the .easy-deploy configuration.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		data, err := os.ReadFile(".easy-deploy")
		if err != nil {
			fmt.Println("unable to read configuration (.easy-deploy):", err)
			os.Exit(1)
		}

		var cfg constants.DeployConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			fmt.Println("invalid config file:", err)
			os.Exit(1)
		}

		if cfg.AwsEbEnv == "" {
			fmt.Println("environment name is not specified in configuration (awsEbEnv)")
			os.Exit(1)
		}

		vars, err := parseKeyValueFile(path)
		if err != nil {
			fmt.Println("failed to parse variables file:", err)
			os.Exit(1)
		}

		ctx := cmd.Context()
		err = aws.UpdateEnvironmentVariables(ctx, cfg.AwsEbEnv, vars)
		if err != nil {
			fmt.Println("failed to update environment:", err)
			os.Exit(1)
		}

		fmt.Println("✅ environment variables pushed to", cfg.AwsEbEnv)
	},
}

func init() {
	rootCmd.AddCommand(pushEnvCmd)
}

func parseKeyValueFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	result := make(map[string]string)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}
		result[parts[0]] = parts[1]
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
