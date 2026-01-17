package main

import (
	"fmt"

	"github.com/zombfeed/GoBlogAggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		return
	}
	err = cfg.SetUser("zomb")
	if err != nil {
		return
	}

	updatedCfg, err := config.Read()
	if err != nil {
		return
	}
	fmt.Printf("%v\n", updatedCfg)
	return
}
