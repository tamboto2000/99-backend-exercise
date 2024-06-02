package snowid

import "github.com/bwmarrin/snowflake"

type Generator interface {
	Generate() int64
}

type SnowID struct {
	node *snowflake.Node
}

func NewSnowID(n int64) (*SnowID, error) {
	node, err := snowflake.NewNode(n)
	if err != nil {
		return nil, err
	}

	snow := SnowID{
		node: node,
	}

	return &snow, nil
}

func (s *SnowID) Generate() int64 {
	return s.node.Generate().Int64()
}

var defSnow Generator

func SetDefault(snow Generator) {
	defSnow = snow
}

func Generate() int64 {
	return defSnow.Generate()
}
