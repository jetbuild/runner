package component

import (
	"fmt"
	"os"
	"plugin"

	"github.com/jetbuild/runner/pkg/component"
)

type Signature func(ctx component.Context) error

func Load(pluginPath string) (Signature, error) {
	if len(pluginPath) == 0 {
		return nil, fmt.Errorf("plugin path does not provided")
	}

	if _, err := os.Stat(pluginPath); err != nil {
		return nil, fmt.Errorf("failed to access component plugin file: %w", err)
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open component plugin: %w", err)
	}

	l, err := p.Lookup("Trigger")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup 'Trigger' plugin symbol: %w", err)
	}

	s, ok := l.(func(ctx component.Context) error)
	if !ok {
		return nil, fmt.Errorf("failed to assert plugin symbol to signature")
	}

	return s, nil
}
