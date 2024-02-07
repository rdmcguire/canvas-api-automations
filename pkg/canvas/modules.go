package canvas

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"k8s.io/utils/ptr"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
)

// TODO this should reeturn error instead of log
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
			log.Error().
				Str("error", err.Error()).
				Msg("Failed listing modules")
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

func (c *Client) GetModuleByID(courseID string, moduleID string) (*canvasauto.Module, error) {
	module := &canvasauto.Module{}
	r, err := c.api.ShowModule(c.ctx, courseID, moduleID, &canvasauto.ShowModuleParams{
		Include: ptr.To([]string{"items"}),
	})
	if err != nil {
		return module, err
	}

	json.NewDecoder(r.Body).Decode(module)
	return module, err
}

func ModuleString(module *canvasauto.Module, showItems bool) string {
	str := fmt.Sprintf("%s [published:%s][id:%s] %d Items",
		StrOrNil(module.Name),
		StrOrNil(module.Published),
		StrOrNil(module.Id),
		len(*module.Items),
	)
	if showItems {
		for _, item := range *module.Items {
			str += "\n\tItem: " + StrOrNil(item.Title)
			str += " Url: " + StrOrNil(item.ExternalUrl)
			str += " ID: " + StrOrNil(item.Id)
		}
	}
	return str
}
