package cmd

import (
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/francescomari/nu/parser"
	"github.com/francescomari/nu/transform"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statsCmd)
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Print statistics about the content",
	Long:  "Reads an export file from stdin and prints statistics about the content on stdout.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		stats, err := transform.Statistics(parser.Parse(os.Stdin))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Computing statistics: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Nodes: %v\n", stats.Nodes)
		fmt.Printf("Properties: %v\n", stats.Properties)
		fmt.Printf("Data: %v\n", size(stats.Data))

		fmt.Printf("Properties per type:\n")
		for _, typ := range sortedStringKeys(stats.PropertiesPerType) {
			fmt.Printf("  %v: %v\n", typ, stats.PropertiesPerType[typ])
		}

		fmt.Printf("Nodes per depth:\n")
		for _, bucket := range sortedIntKeys(stats.NodesPerDepth) {
			fmt.Printf("  %6v: %v\n",
				linearBucket{bucket, transform.StatsNodeDepthBucketSize},
				stats.NodesPerDepth[bucket])
		}

		fmt.Printf("Properties per depth:\n")
		for _, bucket := range sortedIntKeys(stats.PropertiesPerDepth) {
			fmt.Printf("  %6v: %v\n",
				linearBucket{bucket, transform.StatsPropertyDepthBucketSize},
				stats.PropertiesPerDepth[bucket])
		}

		fmt.Printf("Values per size:\n")
		for _, bucket := range sortedIntKeys(stats.ValuesPerSize) {
			fmt.Printf("  %8v: %v\n",
				logarithmicBucket{bucket, transform.StatsValueSizeBucketScale},
				stats.ValuesPerSize[bucket])
		}
	},
}

func sortedStringKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedIntKeys(m map[int]int) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

type humanReadable int64

func (n humanReadable) String() string {
	magnitude := 0
	value := int64(n)

	for value >= 1000 {
		value /= 1000
		magnitude++
	}

	unit := ""

	switch magnitude {
	case 1:
		unit = "kB"
	case 2:
		unit = "MB"
	case 3:
		unit = "GB"
	case 4:
		unit = "PB"
	case 5:
		unit = "EB"
	}

	return fmt.Sprintf("%v%v", value, unit)
}

type size int64

func (n size) String() string {
	if n < 1000 {
		return fmt.Sprintf("%v bytes", int64(n))
	}
	return fmt.Sprintf("%v (%v bytes)", humanReadable(n), int64(n))
}

type linearBucket struct {
	bucket int
	size   int
}

func (b linearBucket) String() string {
	return fmt.Sprintf("%d..%d", b.bucket, b.bucket+b.size)
}

type logarithmicBucket struct {
	bucket int
	base   int
}

func (b logarithmicBucket) String() string {
	return fmt.Sprintf("%v..%v",
		humanReadable(b.begin()),
		humanReadable(b.end()))
}

func (b logarithmicBucket) begin() int64 {
	if b.bucket == 0 {
		return 0
	}
	return b.pow(b.base, b.bucket)
}

func (b logarithmicBucket) end() int64 {
	return b.pow(b.base, b.bucket+1)
}

func (logarithmicBucket) pow(base, exp int) int64 {
	return int64(math.Pow(float64(base), float64(exp)))
}
