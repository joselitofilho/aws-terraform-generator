package config

type RestfulAPI struct {
	Name string `yaml:"name"`
}

func (r *RestfulAPI) GetName() string { return r.Name }
