package main

import (
	"flag"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/netacad"
)

var netacadAssignments []netacad.Assignment
var linkRegexp *regexp.Regexp = regexp.MustCompile(`<a[^>]+href="([^"]+)".*>([^<]+)`)

func assignments() {
	requireArg()
	netacadAssignments = netacad.LoadAssignmentsHtmlFromFile(flag.Args()[1])
	findLinkMatches()
}

// Iterates through all modules, and then in each
// module runs through the list of netacad assignments
// to find matching items within the module.
//
// If a match is found, updateModuleItemLink() will attempt
// to fix any bad links
func findLinkMatches() {
	courseID := flag.Args()[2]
	modules := client.ListModules(courseID)
	for _, module := range modules {
		for _, assignment := range netacadAssignments {
			opts := &canvas.ModuleItemOpts{
				CourseID: courseID,
				Module:   module,
				Name:     assignment.Name,
				URL:      assignment.Url,
				Fuzzy:    false,
			}
			if item := client.GetItemByName(opts); item != nil {
				opts.Item = item
				updateModuleItemLink(opts)
			}
		}
	}
}

// Detects the type of item, then calls the appropriate
// func to adjust bad links
func updateModuleItemLink(opts *canvas.ModuleItemOpts) {
	if opts.Item.ExternalUrl != nil {
		UpdateExternalItemLink(opts)

	} else if canvas.StrStrOrNil(opts.Item.Type) == "Assignment" {
		UpdateAssignmentItemLink(opts)

	} else {
		slog.Error("Unknown Item type", "item", *opts.Item)
	}
}

// Used for items of Type=Assignment
func UpdateAssignmentItemLink(opts *canvas.ModuleItemOpts) {
	aOpts := &canvas.AssignmentOpts{
		ID:             canvas.IntStrOrNil(opts.Item.ContentId),
		ModuleItemOpts: opts,
	}
	assignment, err := client.GetAssignmentById(aOpts)
	if err != nil || assignment == nil {
		slog.Error("Failed to fetch item assignment",
			"assignmentOpts", aOpts,
			"error", err,
			"assignment", assignment,
		)
	}

	// Pull links out of description
	desc := canvas.StrStrOrNil(assignment.Description)
	matches := linkRegexp.FindStringSubmatch(desc)
	if len(matches) != 3 {
		slog.Error("Can't find link in assignment content", "desc", desc)
		return
	}

	// Update each found link in the description
	for _, match := range matches[1:] {
		*assignment.Description = strings.ReplaceAll(desc, match, opts.URL)
	}

	if *assignment.Description == desc {
		slog.Debug("Skipping up-to-date assignment",
			"assignment", canvas.AssignmentString(assignment),
			"inboundDesc", desc)
		return
	}

	slog.Debug("Updating assignment with new description",
		"assignment", canvas.StrStrOrNil(assignment.Name),
		"description", canvas.StrStrOrNil(assignment.Description),
	)
	slog.Info("Updating assignment",
		"assignment", canvas.AssignmentString(assignment))

	// Update the assignment with the new description
	aOpts.Assignment = assignment
	if a, e := client.UpdateAssignment(aOpts); e != nil {
		slog.Error("Failed to update assignment",
			"error", err)
	} else {
		slog.Info("Links updated for assignment",
			"assignment", canvas.AssignmentString(a))
	}
	os.Exit(1)
}

// Used for items that are an external url link
func UpdateExternalItemLink(opts *canvas.ModuleItemOpts) {
	if canvas.StrStrOrNil(opts.Item.ExternalUrl) != opts.URL {
		slog.Warn("Reconciling link for module item",
			"module", canvas.StrStrOrNil(opts.Module.Name),
			"item", canvas.StrStrOrNil(opts.Item.Title),
			"assignment", opts.Name,
			"newLink", opts.URL,
		)
		newItem, err := client.UpdateModuleItemLink(&canvas.ModuleItemOpts{
			CourseID: opts.CourseID,
			Module:   opts.Module,
			Item:     opts.Item,
			Name:     canvas.StrStrOrNil(opts.Item.Title),
			URL:      opts.URL,
		})
		if err != nil {
			panic(err)
		}
		slog.Debug("Item Updated Successfully",
			"item", canvas.ModuleItemString(newItem))
	}
	slog.Debug("Found Match: Module %s, Item %s, Link %s (%s)\n",
		canvas.StrStrOrNil(opts.Module.Name),
		canvas.StrStrOrNil(opts.Item.Title),
		opts.Name,
		opts.URL,
	)
}
