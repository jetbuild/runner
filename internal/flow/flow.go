package flow

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/jetbuild/engine/pkg/flow"
	"github.com/jetbuild/runner/internal/component"
)

type Flow flow.Flow

var components = map[string]component.Signature{}

func (f *Flow) Run() error {
	if err := f.loadComponents(); err != nil {
		return fmt.Errorf("failed to load components: %w", err)
	}

	ctx := component.NewTriggerContext(f.Components[0].Arguments)

	go func() {
		for o := range ctx.OutputStream {
			var outputs sync.Map
			outputs.Store(0, o)
			go f.trigger(&outputs)
		}
	}()

	if err := components[f.Components[0].Key](ctx); err != nil {
		return fmt.Errorf("failed to run trigger component: %w", err)
	}

	return nil
}

func (f *Flow) loadComponents() error {
	d, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	for _, c := range f.Components {
		components[c.Key], err = component.Load(path.Join(d, fmt.Sprintf("%s.so", c.Key)))
		if err != nil {
			return fmt.Errorf("failed to load '%s' component: %w", c.Key, err)
		}
	}

	return nil
}

func (f *Flow) trigger(outputs *sync.Map) {
	targets := f.Components[0].Connections.Targets
	var wg sync.WaitGroup
	wg.Add(len(targets))

	for _, t := range targets {
		if len(targets) == 1 {
			f.runComponent(&wg, outputs, t)

			break
		}

		go f.runComponent(&wg, outputs, t)
	}

	wg.Wait()
}

func (f *Flow) runComponent(wg *sync.WaitGroup, outputs *sync.Map, index uint) {
	defer wg.Done()

	c := f.Components[index]

	if err := components[c.Key](component.NewContext(c.Arguments, outputs, index)); err != nil {
		fmt.Println(err) // TODO
	}

	f.runTargetComponents(wg, outputs, c.Connections.Targets)
}

func (f *Flow) runTargetComponents(wg *sync.WaitGroup, outputs *sync.Map, targets []uint) {
	for _, t := range targets {
		if !f.isComponentSourcesCompleted(outputs, t) {
			break
		}

		wg.Add(1)

		if len(targets) == 1 {
			f.runComponent(wg, outputs, t)

			break
		}

		go f.runComponent(wg, outputs, t)
	}
}

func (f *Flow) isComponentSourcesCompleted(outputs *sync.Map, index uint) bool {
	for _, s := range f.Components[index].Connections.Sources {
		if _, ok := outputs.Load(s); !ok {
			return false
		}
	}

	return true
}
