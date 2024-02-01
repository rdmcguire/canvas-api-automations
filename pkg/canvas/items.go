package canvas

import (
	"encoding/json"
	"fmt"

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
}

func (c *Client) UpdateModuleItemLink(opts *ModuleItemOpts) (*canvasauto.ModuleItem, error) {
	r, err := c.api.UpdateModuleItem(c.ctx,
		opts.CourseID,
		IntStrOrNil(opts.Item.ModuleId),
		IntStrOrNil(opts.Item.Id),
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

func (c *Client) GetItemByName(opts *ModuleItemOpts) *canvasauto.ModuleItem {
	// First try by exact match
	if item := GetItemByTitle(opts.Module.Items, opts.Name); item != nil {
		return item
	}

	// Then get fuzzy if we want to
	if opts.Fuzzy {
		strings := GetItemsStrings(opts.Module.Items)
		matches := fuzzy.FindFold(opts.Name, strings)
		if len(matches) > 1 {
			return GetItemByTitle(opts.Module.Items, matches[0])
		}
	}
	return nil
}

func ModuleItemString(item *canvasauto.ModuleItem) string {
	return fmt.Sprintf("Title: %s, ExternalUrl: %s, ModuleID: %s Published: %s",
		StrStrOrNil(item.Title),
		StrStrOrNil(item.ExternalUrl),
		IntStrOrNil(item.ModuleId),
		BoolStrOrNil(item.Published),
	)
}

func GetItemByTitle(items *[]canvasauto.ModuleItem, title string) *canvasauto.ModuleItem {
	for _, item := range *items {
		if StrStrOrNil(item.Title) == title {
			return &item
		}
	}
	return nil
}

func GetItemsStrings(items *[]canvasauto.ModuleItem) []string {
	strings := make([]string, len(*items))
	for i, item := range *items {
		strings[i] = StrStrOrNil(item.Title)
	}
	return strings
}
