package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nompose",
	Short: "Convert Docker configurations to Nomad jobs",
	Long: `Nompose intelligently converts Docker Compose files, Dockerfiles, 
and Docker images into production-ready Nomad job specifications.

Examples:
  nompose generate docker-compose.yml
  nompose generate Dockerfile  
  nompose generate nginx:latest
  nompose generate ./my-app`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}
