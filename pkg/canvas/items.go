package canvas

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type ModuleItemOpts struct {
	CourseID    string                 // Required for set operations
	Module      *canvasauto.Module     // Module to search
	Item        *canvasauto.ModuleItem // Required for set operations
	Name        string                 // Item title
	URL         string                 // New URL if updating item
	Insensitive bool                   // Enables case-insensitive search
	Fuzzy       bool                   // Enables fuzzy search
	Contains    bool                   // Not so fuzzy but fuzzy-er
}

func (c *Client) UpdateModuleItemLink(opts *ModuleItemOpts) (*canvasauto.ModuleItem, error) {
	r, err := c.api.UpdateModuleItem(c.ctx,
		opts.CourseID,
		StrOrNil(opts.Item.ModuleId),
		StrOrNil(opts.Item.Id),
		canvasauto.UpdateModuleItemJSONRequestBody{
			ModuleItemTitle:       &opts.Name,
			ModuleItemExternalUrl: &opts.URL,
		},
	)
	if err != nil {
		return nil, err
	}

	newItem := new(canvasauto.ModuleItem)
	err = json.NewDecoder(r.Body).Decode(newItem)
	return newItem, err
}

// This version of GetItemByName scans all modules looking for a match
func (c *Client) FindItem(opts *ModuleItemOpts) *canvasauto.ModuleItem {
	for _, module := range c.ListModules(opts.CourseID) {
		opts.Module = module
		if item := c.GetItemByName(opts); item != nil {
			return item
		}
	}
	return nil
}

func (c *Client) GetItemByName(opts *ModuleItemOpts) *canvasauto.ModuleItem {
	// First try by exact match
	if item := GetItemByTitle(opts.Module.Items, opts.Name); item != nil {
		return item
	}

	// Try harder
	itemStrings := GetItemsStrings(opts.Module.Items)
	if opts.Contains {
		for _, s := range itemStrings {
			if strings.Contains(s, opts.Name) {
				return GetItemByTitle(opts.Module.Items, s)
			}
		}
	} else if opts.Fuzzy {
		matches := fuzzy.FindFold(opts.Name, itemStrings)
		if len(matches) > 0 {
			return GetItemByTitle(opts.Module.Items, matches[0])
		}
	}

	return nil
}

func ModuleItemString(item *canvasauto.ModuleItem) string {
	return fmt.Sprintf("Title: %s, ExternalUrl: %s, ModuleID: %s Published: %s",
		StrOrNil(item.Title),
		StrOrNil(item.ExternalUrl),
		StrOrNil(item.ModuleId),
		StrOrNil(item.Published),
	)
}

func GetItemByTitle(items *[]canvasauto.ModuleItem, title string) *canvasauto.ModuleItem {
	for _, item := range *items {
		if StrOrNil(item.Title) == title {
			return &item
		}
	}
	return nil
}

func GetItemsStrings(items *[]canvasauto.ModuleItem) []string {
	strings := make([]string, len(*items))
	for i, item := range *items {
		strings[i] = StrOrNil(item.Title)
	}
	return strings
}
