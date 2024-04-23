package canvas

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
	"k8s.io/utils/ptr"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
)

type ModuleCache struct {
	lock    sync.Mutex
	modules map[int]*canvasauto.Module
}

func (c *ModuleCache) Get(moduleID int) *canvasauto.Module {
	c.lock.Lock()
	defer c.lock.Unlock()
	if module, set := c.modules[moduleID]; set {
		return module
	}
	return nil
}
func (c *ModuleCache) Set(moduleID int, module *canvasauto.Module) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.modules[moduleID] = module
}
func (c *ModuleCache) SetSlice(modules ...*canvasauto.Module) {
	for _, m := range modules {
		c.Set(*m.Id, m)
	}
}

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
			break
		}
		json.NewDecoder(r.Body).Decode(&pageModules)
		modules = append(modules, pageModules...)
		if isLastPage(r) {
			break
		}
		page += 1
	}

	// Populate cache and return
	c.moduleCache.SetSlice(modules...)
	return modules
}

func (c *Client) GetModuleByID(courseID string, moduleID string) (*canvasauto.Module, error) {
	moduleIdInt, _ := strconv.Atoi(moduleID)

	// First hit the cache
	if module := c.moduleCache.Get(moduleIdInt); module != nil {
		return module, nil
	}

	// Then fetch it
	r, err := c.api.ShowModule(c.ctx, courseID, moduleID, &canvasauto.ShowModuleParams{
		Include: ptr.To([]string{"items"}),
	})
	if err != nil {
		return nil, err
	}

	module := &canvasauto.Module{}
	err = json.NewDecoder(r.Body).Decode(module)
	if err != nil {
		return nil, err
	}

	// Add to cache if found
	if module.Name != nil {
		c.moduleCache.Set(moduleIdInt, module)
	}

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
