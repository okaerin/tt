/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/okaerin/tt/internal"
	"github.com/spf13/cobra"
)

var (
	inputFiles []string
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

		// thing := [][]byte{[]byte{1}, []byte{2}}

		for _, i := range inputFiles {
			data, err := os.ReadFile(i)
			if err != nil {
				log.Println(err)
				os.Exit(-1)
			}
			jsons = append(jsons, data)
		}

		// data1 := `{"a":1,"obj":{"c":[1,2,3,4]}}`
		// data2 := `{"b":"1","arr":[{"v":{"o":90}}]}`
		// out, _ := internal.MergeJSONsToJSON([]byte(data1), []byte(data2))
		out, _ := internal.MergeJSONsToJSON(jsons...)
		fmt.Println(string(out))
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
}
