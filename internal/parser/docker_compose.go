package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Jassem-HCP/nompose/internal/types"
	"gopkg.in/yaml.v3"
)

// DockerComposeFile represents the structure of a docker-compose.yml
type DockerComposeFile struct {
	Version  string                                `yaml:"version"`
	Services map[string]types.DockerComposeService `yaml:"services"`
	Networks map[string]interface{}                `yaml:"networks,omitempty"`
	Volumes  map[string]interface{}                `yaml:"volumes,omitempty"`
}

// DockerComposeParser handles parsing docker-compose files
type DockerComposeParser struct{}

// NewDockerComposeParser creates a new parser
func NewDockerComposeParser() *DockerComposeParser {
	return &DockerComposeParser{}
}

// Parse reads and parses a docker-compose file
func (p *DockerComposeParser) Parse(filePath string) ([]types.EnhancedServiceConfig, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read docker-compose file: %w", err)
	}

	// Parse YAML
	var compose DockerComposeFile
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, fmt.Errorf("failed to parse docker-compose YAML: %w", err)
	}

	// Convert to enhanced format
	var services []types.EnhancedServiceConfig
	for name, service := range compose.Services {
		enhanced := types.EnhancedServiceConfig{
			Name:            name,
			OriginalService: service,
			ResolvedImage:   p.getInitialImage(service),
			ResolvedPorts:   p.parsePorts(service.Ports),
			Environment:     p.parseEnvironment(service.Environment),
			Dependencies:    p.parseDependencies(service.DependsOn),
		}
		services = append(services, enhanced)
	}

	return services, nil
}

// getInitialImage determines the initial image (may be placeholder for build)
func (p *DockerComposeParser) getInitialImage(service types.DockerComposeService) string {
	if service.Image != "" {
		return service.Image
	}

	if service.Build != nil {
		buildPath := "."
		switch build := service.Build.(type) {
		case string:
			buildPath = build
		case map[string]interface{}:
			if context, ok := build["context"].(string); ok {
				buildPath = context
			}
		}
		return fmt.Sprintf("{{BUILD_REQUIRED:%s}}", buildPath)
	}

	return "{{NO_IMAGE_SPECIFIED}}"
}

// parsePorts extracts port mappings
func (p *DockerComposeParser) parsePorts(ports []interface{}) []types.PortMapping {
	var mappings []types.PortMapping

	for _, port := range ports {
		switch portVal := port.(type) {
		case string:
			if mapping := p.parsePortString(portVal); mapping != nil {
				mappings = append(mappings, *mapping)
			}
		case int:
			mappings = append(mappings, types.PortMapping{
				Host:      portVal,
				Container: portVal,
				Protocol:  "tcp",
			})
		case map[string]interface{}:
			if mapping := p.parsePortObject(portVal); mapping != nil {
				mappings = append(mappings, *mapping)
			}
		}
	}

	return mappings
}

// parsePortString parses port strings like "8080:80" or "8080"
func (p *DockerComposeParser) parsePortString(portStr string) *types.PortMapping {
	if strings.Contains(portStr, ":") {
		parts := strings.Split(portStr, ":")
		if len(parts) >= 2 {
			host, err1 := strconv.Atoi(parts[0])
			container, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil {
				return &types.PortMapping{
					Host:      host,
					Container: container,
					Protocol:  "tcp",
				}
			}
		}
	} else {
		if port, err := strconv.Atoi(portStr); err == nil {
			return &types.PortMapping{
				Host:      port,
				Container: port,
				Protocol:  "tcp",
			}
		}
	}
	return nil
}

// parsePortObject parses port objects
func (p *DockerComposeParser) parsePortObject(portObj map[string]interface{}) *types.PortMapping {
	mapping := &types.PortMapping{Protocol: "tcp"}

	if target, ok := portObj["target"].(int); ok {
		mapping.Container = target
	}
	if published, ok := portObj["published"].(int); ok {
		mapping.Host = published
	}
	if protocol, ok := portObj["protocol"].(string); ok {
		mapping.Protocol = protocol
	}

	if mapping.Container > 0 && mapping.Host > 0 {
		return mapping
	}
	return nil
}

// parseEnvironment extracts environment variables
func (p *DockerComposeParser) parseEnvironment(env interface{}) map[string]string {
	result := make(map[string]string)

	switch environment := env.(type) {
	case map[string]interface{}:
		for key, value := range environment {
			result[key] = fmt.Sprintf("%v", value)
		}
	case []interface{}:
		for _, item := range environment {
			if envStr, ok := item.(string); ok {
				if strings.Contains(envStr, "=") {
					parts := strings.SplitN(envStr, "=", 2)
					result[parts[0]] = parts[1]
				}
			}
		}
	}

	return result
}

// parseDependencies extracts service dependencies
func (p *DockerComposeParser) parseDependencies(deps interface{}) []string {
	var result []string

	switch depsVal := deps.(type) {
	case []interface{}:
		for _, dep := range depsVal {
			if depStr, ok := dep.(string); ok {
				result = append(result, depStr)
			}
		}
	case map[string]interface{}:
		for serviceName := range depsVal {
			result = append(result, serviceName)
		}
	}

	return result
}
