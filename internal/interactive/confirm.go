package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Jassem-HCP/nompose/internal/types"
)

// Confirmer handles interactive user confirmation and editing
type Confirmer struct {
	scanner *bufio.Scanner
}

// NewConfirmer creates a new interactive confirmer
func NewConfirmer() *Confirmer {
	return &Confirmer{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// ConfirmServices interactively confirms and allows editing of service configurations
func (c *Confirmer) ConfirmServices(services []types.EnhancedServiceConfig) ([]types.EnhancedServiceConfig, error) {
	fmt.Println("\n🔧 Let's review and confirm the detected configurations...")
	fmt.Println("   You can press ENTER to keep detected values, or type new values to override them.")
	fmt.Println()

	var confirmedServices []types.EnhancedServiceConfig

	for i, service := range services {
		fmt.Printf("📦 Service %d/%d: %s\n", i+1, len(services), service.Name)
		fmt.Println(strings.Repeat("─", 50))

		confirmed, err := c.confirmService(service)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm service %s: %w", service.Name, err)
		}

		confirmedServices = append(confirmedServices, confirmed)
		fmt.Println()
	}

	return confirmedServices, nil
}

// confirmService confirms a single service configuration
func (c *Confirmer) confirmService(service types.EnhancedServiceConfig) (types.EnhancedServiceConfig, error) {
	confirmed := service

	// Confirm service name
	if newName, err := c.confirmString("Service name", confirmed.Name, true); err != nil {
		return confirmed, err
	} else if newName != "" {
		confirmed.Name = newName
	}

	// Handle image resolution (build vs image)
	if err := c.resolveImage(&confirmed); err != nil {
		return confirmed, err
	}

	// Confirm ports (show all detected ports)
	if err := c.confirmPorts(&confirmed); err != nil {
		return confirmed, err
	}

	// Show environment variables
	if len(confirmed.Environment) > 0 {
		fmt.Printf("   Environment variables: %d detected ✅\n", len(confirmed.Environment))
		for key, value := range confirmed.Environment {
			fmt.Printf("     %s=%s\n", key, value)
		}
	}

	// Show dependencies
	if len(confirmed.Dependencies) > 0 {
		fmt.Printf("   Dependencies: %v ✅\n", confirmed.Dependencies)
	}

	// Show additional docker-compose settings
	c.showAdditionalSettings(confirmed.OriginalService)

	return confirmed, nil
}

// resolveImage handles image resolution with smart build detection
func (c *Confirmer) resolveImage(service *types.EnhancedServiceConfig) error {
	if strings.HasPrefix(service.ResolvedImage, "{{BUILD_REQUIRED:") {
		// Extract build path
		buildPath := strings.TrimPrefix(service.ResolvedImage, "{{BUILD_REQUIRED:")
		buildPath = strings.TrimSuffix(buildPath, "}}")

		fmt.Printf("   🔨 Build configuration detected\n")
		fmt.Printf("   Build context: %s\n", buildPath)
		fmt.Println()
		fmt.Println("   How would you like to handle the Docker image?")
		fmt.Println("   1. I have the image ready (enter image name/tag)")
		fmt.Println("   2. Build it locally now (auto docker build)")
		fmt.Println("   3. I'll build and push later (enter final image name)")
		fmt.Println()

		choice, err := c.promptForInput("Choice [1-3]", "1", true)
		if err != nil {
			return err
		}

		switch choice {
		case "1":
			imageName, err := c.promptForInput("Enter image name (e.g., my-app:latest)", "", true)
			if err != nil {
				return err
			}
			service.ResolvedImage = imageName

		case "2":
			imageName, err := c.promptForInput("Enter image name to build (e.g., my-app:latest)", fmt.Sprintf("%s:latest", service.Name), true)
			if err != nil {
				return err
			}
			
			fmt.Printf("   🔨 Building image: %s\n", imageName)
			fmt.Printf("   Command: docker build -t %s %s\n", imageName, buildPath)
			fmt.Printf("   ⚠️  Note: Image will be available locally only\n")
			fmt.Printf("   💡 For production, push to registry after building\n")
			
			service.ResolvedImage = imageName

		case "3":
			imageName, err := c.promptForInput("Enter final image name (e.g., registry.com/my-app:latest)", "", true)
			if err != nil {
				return err
			}
			
			fmt.Printf("   📝 You'll need to build and push:\n")
			fmt.Printf("   docker build -t %s %s\n", imageName, buildPath)
			fmt.Printf("   docker push %s\n", imageName)
			
			service.ResolvedImage = imageName

		default:
			service.ResolvedImage = fmt.Sprintf("%s:latest", service.Name)
		}
	} else {
		// Regular image confirmation
		newImage, err := c.confirmString("Image", service.ResolvedImage, true)
		if err != nil {
			return err
		}
		if newImage != "" {
			service.ResolvedImage = newImage
		}
	}
	return nil
}

// confirmPorts handles port confirmation
func (c *Confirmer) confirmPorts(service *types.EnhancedServiceConfig) error {
	if len(service.ResolvedPorts) == 0 {
		fmt.Printf("   Ports: none detected\n")
		return nil
	}

	fmt.Printf("   Ports detected:\n")
	for i, port := range service.ResolvedPorts {
		fmt.Printf("     %d. %d:%d (%s)\n", i+1, port.Host, port.Container, port.Protocol)
	}

	keepPorts, err := c.promptForInput("Keep these port mappings? (Y/n)", "Y", false)
	if err != nil {
		return err
	}

	if strings.ToLower(keepPorts) == "n" || strings.ToLower(keepPorts) == "no" {
		// Allow editing ports (simplified for now)
		fmt.Printf("   Port editing not implemented yet - keeping detected ports\n")
	}

	return nil
}

// showAdditionalSettings displays other docker-compose settings
func (c *Confirmer) showAdditionalSettings(service types.DockerComposeService) {
	if len(service.Volumes) > 0 {
		fmt.Printf("   Volumes: %d detected ✅\n", len(service.Volumes))
		for _, volume := range service.Volumes {
			fmt.Printf("     %s\n", volume)
		}
	}

	if service.WorkingDir != "" {
		fmt.Printf("   Working directory: %s ✅\n", service.WorkingDir)
	}

	if service.User != "" {
		fmt.Printf("   User: %s ✅\n", service.User)
	}

	if service.Restart != "" {
		fmt.Printf("   Restart policy: %s ✅\n", service.Restart)
	}
}

// Helper methods
func (c *Confirmer) confirmString(fieldName, currentValue string, required bool) (string, error) {
	return c.promptForInput(fieldName, currentValue, required)
}

func (c *Confirmer) promptForInput(fieldName, currentValue string, required bool) (string, error) {
	if currentValue != "" {
		fmt.Printf("   %s: %s\n", fieldName, currentValue)
		fmt.Printf("   Keep this value? (Y/n): ")
	} else {
		if required {
			fmt.Printf("   Please enter %s: ", strings.ToLower(fieldName))
		} else {
			fmt.Printf("   Enter %s (optional): ", strings.ToLower(fieldName))
		}
	}

	if !c.scanner.Scan() {
		return "", fmt.Errorf("failed to read user input")
	}

	input := strings.TrimSpace(c.scanner.Text())

	if currentValue != "" {
		if input == "" || strings.ToLower(input) == "y" || strings.ToLower(input) == "yes" {
			return "", nil
		}
		if strings.ToLower(input) == "n" || strings.ToLower(input) == "no" {
			fmt.Printf("   Enter new %s: ", strings.ToLower(fieldName))
			if !c.scanner.Scan() {
				return "", fmt.Errorf("failed to read user input")
			}
			input = strings.TrimSpace(c.scanner.Text())
		}
	}

	if required && input == "" {
		fmt.Printf("   ❌ %s is required. Please enter a value: ", fieldName)
		return c.promptForInput(fieldName, "", required)
	}

	return input, nil
}
