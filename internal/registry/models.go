package registry

import "sort"

type Model struct {
	Provider string
	ModelID  string
}

type RegisteredModel struct {
	Name  string
	Model Model
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

func (r *Registry) List() []RegisteredModel {
	models := make([]RegisteredModel, 0, len(r.models))
	for name, model := range r.models {
		models = append(models, RegisteredModel{
			Name:  name,
			Model: model,
		})
	}

	sort.Slice(models, func(i, j int) bool {
		return models[i].Name < models[j].Name
	})

	return models
}
