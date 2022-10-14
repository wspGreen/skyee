package component

import (
	"reflect"

	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/slog"
	"github.com/wspGreen/skyee/utils"
)

type Component struct {
	iface.IComponent
	IsInit bool
	Name   string
}

func NewComponents() *Components {
	return &Components{
		comps: make(map[string]iface.IComponent),
	}
}

type Components struct {
	comps map[string]iface.IComponent
}

func (cs *Components) AddComponent(c iface.IComponent) {
	name := utils.GetClassName(reflect.TypeOf(c))
	cs.comps[name] = c
	// cs.comps = append(cs.comps, c)
}

func (cs *Components) GetComponent(name string) iface.IComponent {
	c := cs.comps[name]
	return c
}

func (cs *Components) Start() {
	for key, v := range cs.comps {
		slog.Info("Component [%s] Start ", key)
		v.Start()
	}
}
