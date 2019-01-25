package transform

import (
	"encoding/base64"
	"fmt"

	"github.com/francescomari/nu/parser"
)

const (
	// StatsNodeDepthBucketSize is the size of a bucket for the NodesPerDepth
	// field in Stats.
	StatsNodeDepthBucketSize = 10
	// StatsPropertyDepthBucketSize is the size of a bucket for the
	// PropertiesPerDepth field in Stats.
	StatsPropertyDepthBucketSize = 10
	// StatsValueSizeBucketScale is the scale of a bucket for the ValuesPerSize
	// field in Stats.
	StatsValueSizeBucketScale = 1024
)

// Stats contains statistics about the data in an export.
type Stats struct {
	// Nodes is the total amount of nodes in the export.
	Nodes int
	// Properties is the total amount of properties in the export.
	Properties int
	// Data is the total amount of data from every property in the export. For
	// values expressed by a V command, the size is calculated as the length of
	// the value in bytes. For values expressed as an X command, the size is
	// expressed as the length of the Base64-decoded payload.
	Data int64
	// PropertiesPerType is the number of properties grouped by their type.
	PropertiesPerType map[string]int
	// PropertiesPerDepth is the number of properties grouped by the depth of
	// the node they are attached to. PropertiesPerDepth groups the number of
	// properties in buckets of fixed size, where the size of each bucket is
	// StatsPropertyDepthBucketSize.
	PropertiesPerDepth map[int]int
	// NodesPerDepth is the number of nodes grouped by their depth in the
	// content tree. NodesPerDepth groups the number of nodes in buckets of
	// fixed size, where the size of each bucket is StatsNodeDepthBucketSize.
	NodesPerDepth map[int]int
	// ValuesPerSize is the number of property values grouped by their size.
	// ValuesPerSize groups the property values in buckets whose size increases
	// logarithmically. The base of the logarithimc increase is
	// StatsValueSizeBucketScale.
	ValuesPerSize map[int]int
}

// Statistics parses a stream of command and extract statistics about the
// export. Statistics either returns a non-nil Stats or an error.
func Statistics(commands <-chan parser.Cmd) (*Stats, error) {
	stats := Stats{
		PropertiesPerType:  make(map[string]int),
		PropertiesPerDepth: make(map[int]int),
		NodesPerDepth:      make(map[int]int),
		ValuesPerSize:      make(map[int]int),
	}

	if err := stats.parse(commands); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *Stats) parse(commands <-chan parser.Cmd) error {
	for command := range commands {
		switch cmd := command.(type) {
		case parser.R:
			if err := s.parseRoot(commands); err != nil {
				return err
			}
		case parser.Err:
			return s.onError(cmd)
		default:
			return s.onUnexpected(cmd)
		}
	}
	return nil
}

func (s *Stats) parseRoot(commands <-chan parser.Cmd) error {
	s.Nodes++
	s.NodesPerDepth[s.nodeDepthToBucket(0)]++

	for command := range commands {
		switch cmd := command.(type) {
		case parser.C:
			if err := s.parseNode(cmd, 1, commands); err != nil {
				return err
			}
		case parser.P:
			if err := s.parseProperty(cmd, 0, commands); err != nil {
				return err
			}
		case parser.Up:
			return nil
		case parser.Err:
			return s.onError(cmd)
		default:
			return s.onUnexpected(cmd)
		}
	}
	return nil
}

func (s *Stats) parseNode(c parser.C, depth int, commands <-chan parser.Cmd) error {
	s.Nodes++
	s.NodesPerDepth[s.nodeDepthToBucket(depth)]++

	for command := range commands {
		switch cmd := command.(type) {
		case parser.C:
			if err := s.parseNode(cmd, depth+1, commands); err != nil {
				return err
			}
		case parser.P:
			if err := s.parseProperty(cmd, depth, commands); err != nil {
				return err
			}
		case parser.Up:
			return nil
		case parser.Err:
			return s.onError(cmd)
		default:
			return s.onUnexpected(cmd)
		}
	}
	return nil
}

func (s *Stats) parseProperty(p parser.P, depth int, commands <-chan parser.Cmd) error {
	s.Properties++
	s.PropertiesPerType[p.Type]++
	s.PropertiesPerDepth[s.propertyDepthToBucket(depth)]++

	for command := range commands {
		switch cmd := command.(type) {
		case parser.V:
			size := len([]byte(cmd.Data))
			s.Data += int64(size)
			s.ValuesPerSize[s.valueSizeToBucket(size)]++
		case parser.X:
			size := base64.StdEncoding.DecodedLen(len(cmd.Data))
			s.Data += int64(size)
			s.ValuesPerSize[s.valueSizeToBucket(size)]++
		case parser.Up:
			return nil
		case parser.Err:
			return s.onError(cmd)
		default:
			return s.onUnexpected(cmd)
		}
	}
	return nil
}

func (s *Stats) onError(err parser.Err) error {
	return fmt.Errorf("error at line %v: %v", err.Line, err.Err)
}

func (s *Stats) onUnexpected(cmd parser.Cmd) error {
	return fmt.Errorf("unexpected command %T", cmd)
}

func (*Stats) nodeDepthToBucket(depth int) int {
	return (depth / StatsNodeDepthBucketSize) * StatsNodeDepthBucketSize
}

func (*Stats) propertyDepthToBucket(depth int) int {
	return (depth / StatsPropertyDepthBucketSize) * StatsPropertyDepthBucketSize
}

func (*Stats) valueSizeToBucket(size int) int {
	bucket := 0
	for size >= StatsValueSizeBucketScale {
		size = size / StatsValueSizeBucketScale
		bucket++
	}
	return bucket
}
