package cmd

import (
	"errors"
	"fmt"
	"os"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const docsDir = "./docs"

var docsCmd = &cobra.Command{
	Use:     "docs",
	Aliases: []string{"doc", "d"},
	Short:   "Generate documentation for canvas-api-automations",
	Run:     execDocsCmd,
}

func execDocsCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)
	s, e := os.Stat(docsDir)

	docsDir, _ := cmd.Flags().GetString("docsDir")

	// Check if the directory exists, otherwise try to
	// create it
	if errors.Is(e, os.ErrNotExist) {
		if e = os.Mkdir(docsDir, 0755); e != nil {
			log.Fatal().Err(e).Str("docsDir", docsDir).
				Msg("Failed to create docs dir")
		} else {
			s, e = os.Stat(docsDir) // Stat the new dir
		}
	} else if e != nil {
		log.Fatal().Err(e).Msg("Error opening docs dir")
	}

	// Check writable
	if s.Mode().Perm()&0200 != 0200 {
		log.Fatal().
			Str("perms", fmt.Sprintf("%O", s.Mode().Perm())).
			Str("docsDir", docsDir).
			Msg("Docs dir is not writable")
	}

	if err := doc.GenMarkdownTree(rootCmd, docsDir); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate documentation")
	}
}

func init() {
	docsCmd.Flags().String("docsDir", docsDir, "Set output path for generated documentation")
}
