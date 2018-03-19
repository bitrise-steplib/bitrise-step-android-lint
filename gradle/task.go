package gradle

import (
	"fmt"
	"strings"
)

// Task ...
type Task struct {
	name   string
	module Module
}

// GetVariants ...
func (task *Task) GetVariants() (Variants, error) {
	tasksOutput, err := getGradleOutput(task.module.project.location, task.module.name+"tasks")
	if err != nil {
		return nil, fmt.Errorf("%s, %s", tasksOutput, err)
	}

	return cleanStringSlice(task.parseVariants(tasksOutput)), nil
}

func (task *Task) parseVariants(gradleOutput string) Variants {
	tasks := []string{}
lines:
	for _, l := range strings.Split(gradleOutput, "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		l = strings.Split(l, " ")[0]
		if strings.HasPrefix(l, task.name) {
			for _, conflict := range conflicts[task.name] {
				if strings.HasPrefix(l, conflict) {
					continue lines
				}
			}
			l = strings.TrimPrefix(l, task.name)
			if l == "" {
				continue
			}
			tasks = append(tasks, l)
		}
	}
	return tasks
}

// RunVariants ...
func (task *Task) RunVariants(variants Variants) error {
	args := []string{}
	for _, variant := range variants {
		args = append(args, task.module.name+task.name+variant)
	}
	return runGradleCommand(task.module.project.location, args...)
}
