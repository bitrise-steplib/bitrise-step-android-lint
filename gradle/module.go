package gradle

import "fmt"

// Module ...
type Module struct {
	project Project
	name    string
	tasks   []Task
}

// SetTask ...
func (module Module) SetTask(name string) *Task {
	return &Task{
		module: module,
		name:   name,
	}
}

func getGradleModule(configModule string) string {
	if configModule != "" {
		return fmt.Sprintf(":%s:", configModule)
	}
	return ""
}
