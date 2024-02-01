package canvas

import (
	"encoding/json"
	"fmt"
	"strconv"

	"log/slog"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
)

func (c *Client) ListModules(courseID string) []*canvasauto.Module {
	modules := make([]*canvasauto.Module, 0)
	include := []string{"items"}
	page := 1
	for {
		pageModules := make([]*canvasauto.Module, 0)
		pageStr := strconv.Itoa(page)
		r, err := c.api.ListModules(c.ctx,
			courseID,
			&canvasauto.ListModulesParams{
				Page:    &pageStr,
				Include: &include,
			})
		if err != nil {
			slog.Error("Failed listing modules", "error", err)
			continue
		}
		json.NewDecoder(r.Body).Decode(&pageModules)
		modules = append(modules, pageModules...)
		if isLastPage(r) {
			break
		}
		page += 1
	}
	return modules
}

func ModuleString(module *canvasauto.Module) string {
	str := fmt.Sprintf("%s [published:%s][id:%s] %d Items",
		StrStrOrNil(module.Name),
		BoolStrOrNil(module.Published),
		IntStrOrNil(module.Id),
		len(*module.Items),
	)
	for _, item := range *module.Items {
		str += "\n\tItem: " + StrStrOrNil(item.Title)
		str += " Url: " + StrStrOrNil(item.ExternalUrl)
		str += " ID: " + IntStrOrNil(item.Id)
	}
	return str
}
