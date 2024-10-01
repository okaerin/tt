/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/okaerin/tt/internal"
	"github.com/spf13/cobra"
)

var (
	dataFiles     []string
	templateFiles []string
	root          string
	verbose       bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tt -i input.json ... template",
	Short: "template tan",
	Long: `template tan takes multiple JSON files as input and feeds them to a
text/template file.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PreRunE: func(cmd *cobra.Command, args []string) error {
		templateFiles = append(templateFiles, args...)
		if len(templateFiles) < 1 {
			return errors.New("missing template file(s)")
		}
		if root != "" {
			st, err := os.Stat(root)
			if os.IsNotExist(err) {

				return fmt.Errorf(`root dir "%s" not found`, root)
			}
			if !st.Mode().IsDir() {
				return fmt.Errorf(`root "%s"is not a dir `, root)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		lp := internal.NewLogic(verbose)
		lp.Execute(dataFiles, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringSliceVarP(&dataFiles, "input", "i", []string{}, "")
	rootCmd.Flags().StringSliceVarP(&templateFiles, "template", "t", []string{}, "")
	rootCmd.Flags().StringVar(&root, "root", "", "")
	rootCmd.MarkFlagRequired("input")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "")
}
