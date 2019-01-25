package cmd

import (
	"fmt"
	"os"

	"github.com/francescomari/nu/parser"
	"github.com/francescomari/nu/transform"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(propertiesCmd)
}

var propertiesCmd = &cobra.Command{
	Use:   "properties",
	Short: "Print types and fully qualified prooperty paths",
	Long:  "Reads an export file from stdin and prints the type and the fully qualified path of every property on stdout.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for p := range transform.Properties(parser.Parse(os.Stdin)) {
			if p.Err != nil {
				fmt.Fprintf(os.Stderr, "Error at line %v: %v\n", p.Line, p.Err)
				os.Exit(1)
			}
			fmt.Printf("%v %v\n", p.Type, p.Path)
		}
	},
}
