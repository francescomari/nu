package cmd

import (
	"fmt"
	"os"

	"github.com/francescomari/nu/parser"
	"github.com/francescomari/nu/transform"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nodesCmd)
}

var nodesCmd = &cobra.Command{
	Use:   "nodes [path]",
	Short: "Print fully qualified node paths",
	Long:  "Reads an export file from stdin and prints the fully qualified path of every node on stdout.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for p := range transform.Nodes(parser.Parse(os.Stdin)) {
			if p.Err != nil {
				fmt.Fprintf(os.Stderr, "Error at line %v: %v\n", p.Line, p.Err)
				os.Exit(1)
			}
			fmt.Println(p.Path)
		}
	},
}
