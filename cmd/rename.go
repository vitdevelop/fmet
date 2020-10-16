package cmd

import (
	"fmet/utils"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var renameCurrentNames []string
var renameNewNames []string
var renameCurrentRegex string
var renameNewRegex string

func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringArrayVarP(&renameCurrentNames, "current", "c", []string{}, "Current filenames to rename. Note: Use multiple -c.")
	renameCmd.Flags().StringArrayVarP(&renameNewNames, "new", "n", []string{}, "New filenames, same order as current names. Note: Use multiple -n")
	renameCmd.Flags().StringVarP(&renameCurrentRegex, "current-regex", "r", "", "Regex usage in format '^(.+)\\.pdf$/'. Note: Do not use double quotes.")
	renameCmd.Flags().StringVarP(&renameNewRegex, "new-regex", "g", "", "Regex usage in format '$1-new\\.pdf$/'. Note: Do not use double quotes.")
}

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(renameCurrentRegex) > 0 && len(renameNewRegex) > 0 {
			regexRename(renameCurrentRegex, renameNewRegex)
		} else {
			simpleRename(renameCurrentNames, renameNewNames)
		}
	},
}

func renameFile(currentName string, newName string) {
	if !utils.IsPath(currentName) {
		path, err := utils.CurrentPath()
		if err != nil {
			_ = fmt.Errorf("%s\n", err)
		}

		currentName = filepath.Join(path, currentName)
	}
	if !strings.Contains(newName, "/") {
		path, err := os.Getwd()
		if err != nil {
			_ = fmt.Errorf("%s\n", err)
		}
		newName = filepath.Join(path, newName)
	}
	err := os.Rename(currentName, newName)
	if err != nil {
		_ = fmt.Errorf("%s\n", err)
	}
}

func simpleRename(currentNames []string, newNames []string) {
	if len(currentNames) != len(newNames) {
		_ = fmt.Errorf("%s\n", "Current filenames size must be equals with new filenames.")
		os.Exit(1)
	}

	for index, currentName := range currentNames {
		newName := newNames[index]
		renameFile(currentName, newName)
	}
}

func regexRename(currentRegex string, newRegex string) {
	cRegex, err := regexp.Compile(currentRegex)
	if err != nil {
		_ = fmt.Errorf("%s\n", "Current filenames regex is incorrect.")
	}

	if len(workingDirectory) == 0 {
		workingDirectory, err = utils.CurrentPath()
		if err != nil {
			_ = fmt.Errorf("%s\n", "Unknown working directory.")
			os.Exit(1)
		}
	}
	files, err := ioutil.ReadDir(workingDirectory)
	if err != nil {
		_ = fmt.Errorf("%s\n", "Unknown working directory.")
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file != nil && cRegex.MatchString(file.Name()) {
			newFilename := cRegex.ReplaceAllString(file.Name(), newRegex)

			if !dryRun {
				renameFile(filepath.Join(workingDirectory, file.Name()), filepath.Join(workingDirectory, newFilename))
			} else {
				_, _ = fmt.Fprintf(os.Stdout, "%s -> %s\n", file.Name(), newFilename)
			}
		}
	}
}
