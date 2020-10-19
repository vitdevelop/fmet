package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var dryRun bool
var workingDirectory string
var verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run, only print new filenames without affect it. Ex. -d=true")
	rootCmd.PersistentFlags().StringVarP(&workingDirectory, "working-directory", "w", "", "Working directory.")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show info verbose.")
}

var rootCmd = &cobra.Command{
	Use:              "fmet",
	Short:            "fmet is utility for change files metadata",
	Long:             `fmet is a utility for manipulate files metadata and including different formats.`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		if verbose {
			fmt.Printf("Verbose information\n")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}