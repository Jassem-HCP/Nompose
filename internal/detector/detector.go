package detector

import (
	"path/filepath"
	"strings"

	"github.com/Jassem-HCP/nompose/internal/types"
)

// Detector analyzes sources and determines their type
type Detector struct{}

// NewDetector creates a new source detector
func NewDetector() *Detector {
	return &Detector{}
}

// DetectSourceType determines what kind of source we're dealing with
func (d *Detector) DetectSourceType(source string) types.DetectionResult {
	// Clean the source path
	source = strings.TrimSpace(source)
	
	if source == "" {
		return types.DetectionResult{
			Valid: false,
			Error: "empty source provided",
		}
	}

	// GitHub repository detection
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		if strings.Contains(source, "github.com") {
			return types.DetectionResult{
				SourceType: types.SourceGitHubRepo,
				Source:     source,
				Valid:      true,
			}
		}
		return types.DetectionResult{
			Valid: false,
			Error: "URL provided but not a supported GitHub repository",
		}
	}

	// Docker image detection (contains : for tag)
	if strings.Contains(source, ":") && !strings.Contains(source, "/") && !strings.Contains(source, "\\") {
		return types.DetectionResult{
			SourceType: types.SourceDockerImage,
			Source:     source,
			Valid:      true,
		}
	}

	// File-based detection
	filename := filepath.Base(source)
	extension := strings.ToLower(filepath.Ext(source))

	// Docker Compose detection
	if filename == "docker-compose.yml" || 
	   filename == "docker-compose.yaml" ||
	   extension == ".yml" || 
	   extension == ".yaml" {
		return types.DetectionResult{
			SourceType: types.SourceDockerCompose,
			Source:     source,
			Valid:      true,
		}
	}

	// Dockerfile detection
	if strings.HasPrefix(strings.ToLower(filename), "dockerfile") {
		return types.DetectionResult{
			SourceType: types.SourceDockerfile,
			Source:     source,
			Valid:      true,
		}
	}

	// Default to local directory
	return types.DetectionResult{
		SourceType: types.SourceLocalDir,
		Source:     source,
		Valid:      true,
	}
}
