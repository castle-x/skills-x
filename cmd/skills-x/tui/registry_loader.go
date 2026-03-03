package tui

import (
	"fmt"
	"os"
	"sync"

	"github.com/castle-x/skills-x/pkg/registry"
)

var registryWarningsOnce sync.Once

func loadMergedRegistry() (*registry.Registry, error) {
	reg, warnings, err := registry.LoadWithUser()
	if err != nil {
		return nil, err
	}

	registryWarningsOnce.Do(func() {
		for _, w := range warnings {
			fmt.Fprintf(os.Stderr, "⚠ %s\n", w)
		}
	})

	return reg, nil
}
