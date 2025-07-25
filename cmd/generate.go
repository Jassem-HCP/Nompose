package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Jassem-HCP/nompose/internal/detector"
	"github.com/Jassem-HCP/nompose/internal/interactive"
	"github.com/Jassem-HCP/nompose/internal/parser"
)

var generateCmd = &cobra.Command{
	Use:   "generate [source]",
	Short: "Generate Nomad jobs from Docker configurations",
	Long: `Analyze Docker configurations and generate Nomad job specifications.

Supported sources:
  - docker-compose.yml (→ multiple Nomad jobs)
  - Dockerfile (→ single Nomad job) 
  - nginx:latest (→ single Nomad job from image)
  - ./my-app (→ analyze directory)
  - https://github.com/user/repo (→ analyze repository)`,
	Example: `  nompose generate docker-compose.yml
  nompose generate Dockerfile
  nompose generate nginx:alpine
  nompose generate ./my-project`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	source := args[0]
	
	fmt.Printf("🔍 Analyzing source: %s\n", source)
	
	// Detect source type
	detector := detector.NewDetector()
	result := detector.DetectSourceType(source)
	
	if !result.Valid {
		return fmt.Errorf("❌ %s", result.Error)
	}
	
	fmt.Printf("✅ Detected source type: %s\n", result.SourceType)
	
	// Parse based on source type
	switch result.SourceType {
	case "docker-compose":
		return handleDockerCompose(source)
	case "dockerfile":
		return handleDockerfile(source)
	case "docker-image":
		return handleDockerImage(source)
	default:
		fmt.Printf("🚧 %s parsing coming in next sub-steps...\n", result.SourceType)
	}
	
	return nil
}

func handleDockerCompose(filePath string) error {
	fmt.Printf("📋 Parsing docker-compose file...\n")
	
	// Parse with enhanced data preservation
	parser := parser.NewDockerComposeParser()
	services, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse docker-compose: %w", err)
	}
	
	fmt.Printf("✅ Found %d services:\n", len(services))
	
	// Show enhanced detection summary
	for i, service := range services {
		fmt.Printf("   %d. %s (%s)\n", i+1, service.Name, service.ResolvedImage)
		if len(service.ResolvedPorts) > 0 {
			fmt.Printf("      Ports: %d detected\n", len(service.ResolvedPorts))
		}
		if len(service.Environment) > 0 {
			fmt.Printf("      Environment: %d variables\n", len(service.Environment))
		}
	}
	
	// Enhanced interactive confirmation
	confirmer := interactive.NewConfirmer()
	confirmedServices, err := confirmer.ConfirmServices(services)
	if err != nil {
		return fmt.Errorf("failed to confirm services: %w", err)
	}
	
	// Show final summary
	fmt.Println("✅ Configuration confirmed! Summary:")
	for i, service := range confirmedServices {
		fmt.Printf("📦 %d. %s\n", i+1, service.Name)
		fmt.Printf("   Image: %s\n", service.ResolvedImage)
		if len(service.ResolvedPorts) > 0 {
			fmt.Printf("   Ports: %d configured\n", len(service.ResolvedPorts))
		}
		fmt.Printf("   Environment: %d variables\n", len(service.Environment))
		if len(service.Dependencies) > 0 {
			fmt.Printf("   Dependencies: %v\n", service.Dependencies)
		}
		fmt.Println()
	}
	
	fmt.Println("🚧 Enhanced Nomad job generation coming next...")
	
	return nil
}

func handleDockerfile(filePath string) error {
	fmt.Printf("🚧 Dockerfile parsing coming in next sub-step...\n")
	return nil
}

func handleDockerImage(image string) error {
	fmt.Printf("🚧 Docker image analysis coming in next sub-step...\n")
	return nil
}
