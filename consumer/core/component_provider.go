package core

import (
	"github.com/golobby/container/v3"
	"github.com/latifrons/goerrorcollector/consumer/service"
	"github.com/latifrons/latigo/program"
)

type ComponentProvider struct {
	DisabledComponents map[string]bool
}

func (c *ComponentProvider) ProvideAllComponents() []program.Component {

	var collector *service.ErrorListener
	err := container.Resolve(&collector)
	if err != nil {
		panic(err)
	}

	return []program.Component{
		collector,
	}
}

func (c *ComponentProvider) ProvideDisabledComponents() map[string]bool {
	return c.DisabledComponents
}
