package infra

import "github.com/tamboto2000/99-backend-exercise/pkg/snowid"

func InitSnowID() error {
	// TODO: get nodeId from config
	nodeId := 1
	node, err := snowid.NewSnowID(int64(nodeId))
	if err != nil {
		return err
	}

	snowid.SetDefault(node)

	return nil
}
