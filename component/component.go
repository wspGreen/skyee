package component

type IComponent interface {
	Start()
	Stop()
}

type Component struct {
	IComponent
	IsInit bool
}

type Components struct {
	comps []IComponent
}

func (cs *Components) AddComponent(c IComponent) {
	cs.comps = append(cs.comps, c)
}

func (cs *Components) Start() {
	for _, v := range cs.comps {
		v.Start()
	}
}
