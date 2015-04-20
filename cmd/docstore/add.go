package docstore

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/andrew-d/docstore/store"
)

var (
	addDocumentCmd = &cobra.Command{
		Use:     "document <document name> <file>...",
		Aliases: []string{"doc"},
		Short:   "Add a new document (and optionally some files) to the store",
		Run:     wrapCommand(runAddDocument),
	}

	flagDocumentTags StringSlice
)

func init() {
	addDocumentCmd.Flags().VarP(&flagDocumentTags, "tag", "t",
		"Tags to add to the new document")

	addCmd.AddCommand(addDocumentCmd)
}

func runAddDocument(cmd *cobra.Command, args []string, store *store.Store) {
	if len(args) < 1 {
		log.Error("No document name provided")
		cmd.Help()
		return
	}

	files := []*os.File{}

	// Open files to verify they're all there
	for _, fname := range args[1:] {
		f, err := os.Open(fname)
		if err != nil {
			log.WithFields(logrus.Fields{
				"err":      err,
				"filename": fname,
			}).Error("Could not open file")
			return
		}
		defer f.Close()

		files = append(files, f)
	}

	// Create new document
	doc, err := store.CreateDocument(args[0])
	if err != nil {
		log.WithField("err", err).Error("Could not create document")
		return
	}
	log.WithField("id", doc.ID).Debug("Created document")

	// Add files to this document
	for i, f := range files {
		log.WithField("filename", args[1+i]).Debug("Adding file to store...")

		// TODO
		_ = doc
		_ = f

		log.WithField("filename", args[1+i]).Debug("File added")
	}

	// TODO: add tags to the document from flagDocumentTags
	log.WithField("tags", flagDocumentTags).Debug("Adding tags to document")

	// TODO: what do we return?
}
