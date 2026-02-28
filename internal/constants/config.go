package constants

type DeployConfig struct {
	Framework          string `json:"framework"`
	DeploymentSolution string `json:"deploymentSolution"`
	Branch             string `json:"branch"`
	LaunchInVPC        bool   `json:"launchInVpc"`

	AwsAccessKeyID     string `json:"awsAccessKeyId,omitempty"`
	AwsSecretAccessKey string `json:"awsSecretAccessKey,omitempty"`
	AwsRegion          string `json:"awsRegion,omitempty"`
	AwsS3Bucket        string `json:"awsS3Bucket,omitempty"`
	AwsEbApp           string `json:"awsEbApp,omitempty"`
	AwsEbEnv           string `json:"awsEbEnv,omitempty"`
}

var FrameworkMap = map[string]string{
	"Node.js": "nodejs",
	"Python":  "python",
	"Java":    "java",
	".NET":    "dotnet",
	"Ruby":    "ruby",
}

var DeploymentMap = map[string]string{
	"GitHub": "github",
	"GitLab": "gitlab",
}

var SolutionStackKeywordMap = map[string]string{
	"nodejs": "Node.js 22",
	"python": "Python 3.11",
	"java":   "Corretto 17",
	"dotnet": ".NET 8",
	"ruby":   "Ruby 3.3",
}

var SolutionStackFallbackMap = map[string]string{
	"nodejs": "64bit Amazon Linux 2023 v6.8.0 running Node.js 22",
	"python": "64bit Amazon Linux 2023 v4.4.0 running Python 3.11",
	"java":   "64bit Amazon Linux 2023 v5.3.0 running Corretto 17",
	"dotnet": "64bit Amazon Linux 2023 v4.1.0 running .NET 8",
	"ruby":   "64bit Amazon Linux 2023 v4.1.0 running Ruby 3.3",
}
