package cmd

import (
	"fmt"
	"os"

	"github.com/francescomari/nu/filter"
	"github.com/francescomari/nu/parser"
	"github.com/francescomari/nu/serializer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(subtreeCmd)
}

var subtreeCmd = &cobra.Command{
	Use:   "subtree",
	Short: "Shrinks the export to a subtree",
	Long:  "Reads an export file from stdin, shrinks it to a specific subtree, and prints the resulting export on stdout.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := filter.Subtree(args[0], parser.Parse(os.Stdin))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid argument: %v\n", err)
			os.Exit(1)
		}
		if err := serializer.Serialize(out, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Error while serializing: %v\n", err)
			os.Exit(1)
		}
	},
}
