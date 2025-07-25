package types

// SourceType represents different input sources
type SourceType string

const (
	SourceDockerCompose SourceType = "docker-compose"
	SourceDockerfile    SourceType = "dockerfile"  
	SourceDockerImage   SourceType = "docker-image"
	SourceLocalDir      SourceType = "local-directory"
	SourceGitHubRepo    SourceType = "github-repo"
)

// DetectionResult holds what we detected about the source
type DetectionResult struct {
	SourceType SourceType
	Source     string
	Valid      bool
	Error      string
}

// EnhancedServiceConfig preserves all docker-compose data
type EnhancedServiceConfig struct {
	Name            string
	OriginalService DockerComposeService
	ResolvedImage   string                 // Final image after user input
	ResolvedPorts   []PortMapping         // Processed port mappings
	Environment     map[string]string      // Flattened environment
	Dependencies    []string               // Flattened dependencies
}

// PortMapping represents a port configuration
type PortMapping struct {
	Host      int    // Host port (8080)
	Container int    // Container port (80)
	Protocol  string // tcp/udp
}

// DockerComposeService mirrors the docker-compose service structure
type DockerComposeService struct {
	Image       string                 `yaml:"image,omitempty"`
	Build       interface{}            `yaml:"build,omitempty"`
	Ports       []interface{}          `yaml:"ports,omitempty"`
	Environment interface{}            `yaml:"environment,omitempty"`
	Volumes     []string               `yaml:"volumes,omitempty"`
	DependsOn   interface{}            `yaml:"depends_on,omitempty"`
	HealthCheck *HealthCheckConfig     `yaml:"healthcheck,omitempty"`
	Deploy      *DeployConfig          `yaml:"deploy,omitempty"`
	Restart     string                 `yaml:"restart,omitempty"`
	Networks    interface{}            `yaml:"networks,omitempty"`
	Command     interface{}            `yaml:"command,omitempty"`
	WorkingDir  string                 `yaml:"working_dir,omitempty"`
	User        string                 `yaml:"user,omitempty"`
	Labels      map[string]string      `yaml:"labels,omitempty"`
	Expose      []string               `yaml:"expose,omitempty"`
}

// HealthCheckConfig represents healthcheck configuration
type HealthCheckConfig struct {
	Test        interface{} `yaml:"test,omitempty"`
	Interval    string      `yaml:"interval,omitempty"`
	Timeout     string      `yaml:"timeout,omitempty"`
	Retries     int         `yaml:"retries,omitempty"`
	StartPeriod string      `yaml:"start_period,omitempty"`
}

// DeployConfig represents deploy configuration  
type DeployConfig struct {
	Replicas int `yaml:"replicas,omitempty"`
}

// Legacy types for backward compatibility (if needed)
type ServiceConfig struct {
	Name        string
	Image       string
	Port        int
	Environment map[string]string
	Volumes     []string
	DependsOn   []string
	HealthCheck string
	Replicas    int
}

type GenerateOptions struct {
	ServiceName  string
	Port         int
	Instances    int
	CPU          int
	Memory       int
	HealthCheck  string
	Datacenter   string
	Namespace    string
	OutputFile   string
	OutputFormat string
	DryRun       bool
	Interactive  bool
	ForceType    string
	WithConsul   bool
	WithVault    bool
	WithIngress  bool
}
