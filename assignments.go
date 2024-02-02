package main

import (
	"flag"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/netacad"
)

// Set this to false if an assignnment containing
// the string "Final Exam" should be used. This is due to
// the option for "Final Comprehensive Exam" and the awful job
// whomever manually click-grunted these assignments into Canvas did
const skipFinalNotComprehensive = true

var (
	linkRegexp *regexp.Regexp = regexp.MustCompile(`<a[^>]+href="([^"]+)".*>([^<]+)`)
	examRegexp *regexp.Regexp = regexp.MustCompile(`Chapter (\d+) Exam`)
	znumRegexp *regexp.Regexp = regexp.MustCompile(` 0(\d+)`)
)

var netacadAssignments []netacad.Assignment

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

	for _, assignment := range netacadAssignments {
		// Whatever. It is what it is.
		if strings.Contains(assignment.Name, "Final E") && skipFinalNotComprehensive {
			log.Warn().Msg("Skipping non-comprehensive final. Set const skipFinalNotComprehensive=true to force")
			continue
		}
		opts := &canvas.ModuleItemOpts{
			CourseID: courseID,
			Name:     assignment.Name,
			URL:      assignment.Url,
		}

	Assignment:
		for _, module := range modules {
			opts.Module = module
			opts.Item = findAssignmentInModule(opts)

			if opts.Item != nil {
				updateModuleItemLink(opts)
				break Assignment
			}
		}

		if opts.Item == nil {
			log.Info().Any("assignment", assignment).
				Msg("Netacad assignment not found in Canvas. Should it be?")
		}
	}
}

func findAssignmentInModule(opts *canvas.ModuleItemOpts) *canvasauto.ModuleItem {
	// Try an exact match
	if item := client.GetItemByName(opts); item != nil {
		return item
	}

	origName := opts.Name

	// Try to fix leading 0's
	opts.Name = znumRegexp.ReplaceAllString(opts.Name, " $1")
	if item := client.GetItemByName(opts); item != nil {
		log.Debug().
			Str("originalName", origName).
			Str("foundItem", canvas.StrStrOrNil(item.Title)).
			Msg("Replace leading zero match")
		return item
	}
	opts.Name = origName

	// Try to fix Exam -> Quiz
	if strings.Contains(opts.Name, "Exam") {
		opts.Name = rewriteExamToQuiz(opts.Name)
		if item := client.GetItemByName(opts); item != nil {
			log.Debug().
				Str("originalName", origName).
				Str("foundItem", canvas.StrStrOrNil(item.Title)).
				Msg("Rewrite Exam to Quiz Result")
			return item
		}
	}
	opts.Name = origName

	// Try to fix midterm/final chapter vs module stupidity
	if strings.Contains(opts.Name, " (M") {
		opts.Name = strings.Split(opts.Name, "(")[0]
		opts.Name = strings.Trim(opts.Name, " ")
		if item := client.GetItemByName(opts); item != nil {
			log.Debug().
				Str("originalName", origName).
				Str("foundItem", canvas.StrStrOrNil(item.Title)).
				Msg("Rewrite Module Parenthesis")
			return item
		}
	}

	// Lastly, be fuzzy
	opts.Fuzzy = true
	if item := client.GetItemByName(opts); item != nil {
		log.Info().Str("foundItem", canvas.StrStrOrNil(item.Title)).
			Msg("Fuzzy result")
	}

	opts.Name = origName
	return nil
}

func rewriteExamToQuiz(name string) string {
	match := examRegexp.FindStringSubmatch(name)
	if len(match) == 2 {
		return "Quiz " + strings.TrimPrefix(match[1], "0")
	}
	return name
}

// Detects the type of item, then calls the appropriate
// func to adjust bad links
func updateModuleItemLink(opts *canvas.ModuleItemOpts) {
	if opts.Item.ExternalUrl != nil {
		UpdateExternalItemLink(opts)

	} else if canvas.StrStrOrNil(opts.Item.Type) == "Assignment" {
		UpdateAssignmentItemLink(opts)

	} else {
		log.Error().Any("tem", *opts.Item).Msg("Unsupported item type")
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
		log.Error().
			Any("assignmentOpts", aOpts).
			Any("error", err).
			Any("assignment", assignment).
			Msg("Failed to fetch item assignment")
	}

	// Pull links out of description
	desc := canvas.StrStrOrNil(assignment.Description)
	matches := linkRegexp.FindStringSubmatch(desc)
	if len(matches) != 3 {
		log.Info().Str("desc", desc).
			Msg("Can't find link in assignment content")
		return
	}

	// Update each found link in the description
	for _, match := range matches[1:] {
		*assignment.Description = strings.ReplaceAll(desc, match, opts.URL)
	}

	if *assignment.Description == desc {
		log.Debug().
			Str("assignment", canvas.AssignmentString(assignment)).
			Msg("Skipping up-to-date assignment")
		return
	}

	log.Debug().
		Str("assignment", canvas.StrStrOrNil(assignment.Name)).
		Str("description", canvas.StrStrOrNil(assignment.Description)).
		Msg("Updating assignment with new description")
	log.Info().
		Str("assignment", canvas.AssignmentString(assignment)).
		Msg("Updating assignment")

	// Update the assignment with the new description
	aOpts.Assignment = assignment
	if a, e := client.UpdateAssignment(aOpts); e != nil {
		log.Error().
			Str("error", err.Error()).
			Msg("Failed to update assignment")
	} else {
		log.Info().
			Str("assignment", canvas.AssignmentString(a)).
			Msg("Links updated for assignment")
	}
}

// Used for items that are an external url link
func UpdateExternalItemLink(opts *canvas.ModuleItemOpts) {
	if canvas.StrStrOrNil(opts.Item.ExternalUrl) != opts.URL {
		log.Warn().
			Str("module", canvas.StrStrOrNil(opts.Module.Name)).
			Str("item", canvas.StrStrOrNil(opts.Item.Title)).
			Str("assignment", opts.Name).
			Str("newLink", opts.URL).
			Msg("Reconciling link for module item")
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
		log.Debug().
			Str("item", canvas.ModuleItemString(newItem)).
			Msg("Item Updated Successfully")
	}
	log.Debug().
		Str("Module", canvas.StrStrOrNil(opts.Module.Name)).
		Str("Item", canvas.StrStrOrNil(opts.Item.Title)).
		Str("Title", opts.Name).
		Str("Link", opts.URL).
		Msg("Found Match")
}
