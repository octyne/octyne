package registry

type Model struct {
	Provider string
	ModelID  string
}

type Registry struct {
	models map[string]Model
}

func NewRegistry() *Registry {
	return &Registry{
		models: make(map[string]Model),
	}
}

func (r *Registry) Register(name string, model Model) {
	r.models[name] = model
}

func (r *Registry) Get(name string) (Model, bool) {
	model, ok := r.models[name]
	return model, ok
}
