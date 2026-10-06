// Package model contains shared declarative values, without execution or filesystem dependencies.
package model

type Command struct {
	Args        []string `yaml:"args,omitempty" json:"args,omitempty"`
	Shell       string   `yaml:"shell,omitempty" json:"shell,omitempty"`
	Group       string   `yaml:"group,omitempty" json:"group,omitempty"`
	Destructive bool     `yaml:"destructive,omitempty" json:"destructive,omitempty"`
	Requires    []string `yaml:"requires,omitempty" json:"requires,omitempty"`
}
type Infrastructure struct {
	Kind      string `yaml:"kind,omitempty" json:"kind,omitempty"`
	File      string `yaml:"file,omitempty" json:"file,omitempty"`
	Mode      string `yaml:"mode,omitempty" json:"mode,omitempty"`
	Context   string `yaml:"context,omitempty" json:"context,omitempty"`
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Image     string `yaml:"image,omitempty" json:"image,omitempty"`
	Port      int    `yaml:"port,omitempty" json:"port,omitempty"`
	Generate  bool   `yaml:"generate,omitempty" json:"generate,omitempty"`
}
type Database struct {
	Kind     string `yaml:"kind,omitempty" json:"kind,omitempty"`
	Tool     string `yaml:"tool,omitempty" json:"tool,omitempty"`
	Generate bool   `yaml:"generate,omitempty" json:"generate,omitempty"`
	URLenv   string `yaml:"url_env,omitempty" json:"url_env,omitempty"`
	Path     string `yaml:"path,omitempty" json:"path,omitempty"`
}
type Service struct {
	ID             string             `yaml:"id" json:"id"`
	Dir            string             `yaml:"dir" json:"dir"`
	Language       string             `yaml:"language,omitempty" json:"language,omitempty"`
	Manager        string             `yaml:"manager,omitempty" json:"manager,omitempty"`
	GoWorkspace    string             `yaml:"go_workspace,omitempty" json:"go_workspace,omitempty"`
	Framework      string             `yaml:"framework,omitempty" json:"framework,omitempty"`
	Infrastructure Infrastructure     `yaml:"infrastructure" json:"infrastructure"`
	Database       Database           `yaml:"database,omitempty" json:"database,omitempty"`
	Commands       map[string]Command `yaml:"commands,omitempty" json:"commands,omitempty"`
	Pending        []string           `yaml:"pending,omitempty" json:"pending,omitempty"`
}
type Suggestion struct {
	Purpose  string  `json:"purpose,omitempty"`
	Service  string  `json:"service"`
	Name     string  `json:"name"`
	Command  Command `json:"command"`
	Evidence string  `json:"evidence"`
}
type Action struct {
	Name, Group, Blocked, Special string
	Command                       Command
}
