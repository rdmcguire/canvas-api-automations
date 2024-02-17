package util

import (
	"fmt"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

func MustFuzzyFindModule(cmd *cobra.Command) *canvasauto.Module {
	var module *canvasauto.Module
	var err error
	if module, err = FuzzyFindModule(cmd); err != nil {
		Logger(cmd).Fatal().Err(err).Msg("Failed to select module")
	}
	return module
}

func FuzzyFindModule(cmd *cobra.Command) (*canvasauto.Module, error) {
	client := Client(cmd)

	modules := client.ListModules(GetCourseIdStr(cmd))
	idx, err := fuzzyfinder.Find(modules, func(i int) string {
		return fmt.Sprintf("%s - %s",
			canvas.StrOrNil(modules[i].Id),
			canvas.StrOrNil(modules[i].Name))
	})

	if err != nil {
		return nil, err
	}
	return modules[idx], nil
}

func MustFuzzyFindAssignment(cmd *cobra.Command, module *canvasauto.Module) *canvasauto.Assignment {
	var assignment *canvasauto.Assignment
	var err error
	if assignment, err = FuzzyFindAssignment(cmd, module); err != nil {
		Logger(cmd).Fatal().Err(err).Msg("Failed to select module assignment")
	}
	return assignment
}

func FuzzyFindAssignment(cmd *cobra.Command, module *canvasauto.Module) (*canvasauto.Assignment, error) {
	log := Logger(cmd)
	client := Client(cmd)

	assignments := client.GetAssignmentsFromModule(GetCourseIdStr(cmd), module)
	if len(assignments) < 1 {
		log.Fatal().Str("id", canvas.StrOrNil(module.Id)).
			Str("name", canvas.StrOrNil(module.Name)).
			Msg("No assignments available for module")
	}

	idx, err := fuzzyfinder.Find(assignments, func(i int) string {
		return fmt.Sprintf("%s - %s - %s",
			canvas.StrOrNil(module.Name),
			canvas.StrOrNil(assignments[i].Id),
			canvas.StrOrNil(assignments[i].Name))
	})
	if err != nil {
		return nil, err
	}

	return assignments[idx], nil
}
