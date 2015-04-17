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
)

func init() {
	showCmd.AddCommand(showDocumentCmd)
}

func runShowDocument(cmd *cobra.Command, args []string, store *store.Store) {
	if len(args) != 1 {
		log.Error("Command takes exactly one argument")
		cmd.Help()
		return
	}

	// Is this a number or a string?
	var (
		id    uint64
		d     models.Document
		found bool
		err   error
	)

	if id, err = strconv.ParseUint(args[0], 10, 64); err != nil {
		d, found, err = store.GetDocumentById(int64(id))
	} else {
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

	// TODO: print in a real format?
	json.NewEncoder(os.Stdout).Encode(&d)
}
