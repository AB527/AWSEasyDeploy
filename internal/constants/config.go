package constants

// DeployConfig stores normalized values
type DeployConfig struct {
	Framework          string `json:"framework"`
	DeploymentSolution string `json:"deployment_solution"`
	Branch             string `json:"branch"`
	LaunchInVPC        bool   `json:"launch_in_vpc"`
}

// Map input to normalized values
var FrameworkMap = map[string]string{
	"Go":      "go",
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
