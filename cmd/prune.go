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
	rootCmd.AddCommand(pruneCmd)
}

var pruneCmd = &cobra.Command{
	Use:   "prune [path]",
	Short: "Remove a subtree from an export",
	Long:  "Reads an export file from stdin, remove a subtree from it, and prints the resulting export on stdout.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := filter.Prune(args[0], parser.Parse(os.Stdin))
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
