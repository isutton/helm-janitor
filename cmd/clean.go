package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cleanCmd *cobra.Command

var dryRun bool

func init() {
	// cleanCmd represents the clean command
	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Removes artifacts of previous releases",
		Run: func(cmd *cobra.Command, args []string) {
			cleanUp()
		},
	}

	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run")
}

func cleanUp() {
	if (dryRun) {
		fmt.Println("Dry run called...")
	} else {
		fmt.Println("The real thing called...")
	}
}
