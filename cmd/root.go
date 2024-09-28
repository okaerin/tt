/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"text/template"

	"github.com/okaerin/tt/internal"
	"github.com/spf13/cobra"
)

var (
	inputFiles []string
	verbose    bool
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
		if len(args) < 1 {
			return errors.New("missing template file(s)")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		var jsons [][]byte

		for _, i := range inputFiles {
			data, err := os.ReadFile(i)
			if err != nil {
				log.Println(err)
				os.Exit(-1)
			}
			jsons = append(jsons, data)
		}

		mergedJSONs, _ := internal.MergeJSONsToJSON(jsons...)
		//parse back to map
		mergedMap := internal.JSONToMap(mergedJSONs)

		if verbose {
			log.Println(string(mergedJSONs))
		}

		for _, arg := range args {
			data, err := os.ReadFile(arg)
			if err != nil {
				log.Printf(`could not parse template "%s"`, arg)
			}
			tmpl := template.Must(template.New("").Funcs(template.FuncMap{
				"toJSON": func(data interface{}) string {
					if jstr, err := json.Marshal(data); err == nil {
						return string(jstr)
					}
					return ""
				},
			}).Parse(string(data)))

			tmpl.Execute(os.Stdout, mergedMap)
		}
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
	rootCmd.Flags().StringSliceVarP(&inputFiles, "input", "i", []string{}, "")
	rootCmd.MarkFlagRequired("input")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "")
}
