/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show semantic version adn release date.",
	Long:  "Show semantic version adn release date.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s (releasedAt: %s)\n", version, releaseDate)
	},
}

func init() {
}
