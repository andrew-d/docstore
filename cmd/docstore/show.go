package docstore

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/andrew-d/docstore/store"
	"github.com/andrew-d/docstore/store/models"
)

var (
	showDocumentCmd = &cobra.Command{
		Use:     "document <id or name>",
		Aliases: []string{"doc"},
		Short:   "Show a single document, given either an ID or a name",
		Run:     wrapCommand(runShowDocument),
	}

	showTagCmd = &cobra.Command{
		Use:   "tag <id or name>",
		Short: "Show a single tag, given either an ID or a name",
		Run:   wrapCommand(runShowTag),
	}

	flagShowFormat string
)

func init() {
	showDocumentCmd.Flags().StringVar(&flagShowFormat, "format", "",
		"output format to use - will use default if not set")

	showCmd.AddCommand(showDocumentCmd)
	showCmd.AddCommand(showTagCmd)
}

func runShowDocument(cmd *cobra.Command, args []string, store *store.Store) {
	if len(args) != 1 {
		log.Error("Command takes exactly one argument")
		cmd.Help()
		return
	}

	if len(flagShowFormat) == 0 {
		flagShowFormat = "Document {{.ID}} - {{.Name}} - {{len .Files}} file(s)\n"
	}

	tmpl, err := ParseFormat(flagShowFormat)
	if err != nil {
		log.WithField("err", err).Error("Invalid format")
		return
	}

	// Is this a number or a string?
	var (
		id    uint64
		d     models.Document
		found bool
	)

	if id, err = strconv.ParseUint(args[0], 10, 64); err == nil {
		log.Debug("Document specifier is numeric")
		d, found, err = store.GetDocumentById(int64(id))
	} else {
		log.Debug("Document specifier is a string")
		d, found, err = store.GetDocumentByName(args[0])
	}

	if err != nil {
		log.WithField("err", err).Error("Error fetching document")
		return
	}
	if !found {
		log.Info("Document not found")
		return
	}

	err = tmpl.Execute(os.Stdout, d)
	if err != nil {
		log.WithField("err", err).Error("Error rendering show template")
	}
}

func runShowTag(cmd *cobra.Command, args []string, store *store.Store) {
	if len(args) != 1 {
		log.Error("Command takes exactly one argument")
		cmd.Help()
		return
	}

	if len(flagShowFormat) == 0 {
		// TODO: number of documents
		flagShowFormat = "Tag {{.ID}} - {{.Name}\n"
	}

	tmpl, err := ParseFormat(flagShowFormat)
	if err != nil {
		log.WithField("err", err).Error("Invalid format")
		return
	}

	// Is this a number or a string?
	var (
		id    uint64
		t     models.Tag
		found bool
	)

	if id, err = strconv.ParseUint(args[0], 10, 64); err != nil {
		t, found, err = store.GetTagById(int64(id))
	} else {
		t, found, err = store.GetTagByName(args[0])
	}

	if err != nil {
		log.WithField("err", err).Error("Error fetching tag")
		return
	}
	if !found {
		log.Info("Tag not found")
		return
	}

	err = tmpl.Execute(os.Stdout, d)
	if err != nil {
		log.WithField("err", err).Error("Error rendering show template")
	}
}
