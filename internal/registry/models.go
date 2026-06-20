package registry

type Model struct {
	Provider string
	ModelID  string
}

var models = map[string]Model{
	"gpt-4.1-mini": {
		Provider: "openai",
		ModelID:  "gpt-4.1-mini",
	},
}

func Get(name string) (Model, bool) {
	model, ok := models[name]
	return model, ok
}
