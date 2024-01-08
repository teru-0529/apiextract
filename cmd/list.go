/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents output url and httpmethod list
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Output url and httpmethod list.",
	Long:  "Output url and httpmethod list.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("list called")
		return nil
	},
}

func init() {
}
