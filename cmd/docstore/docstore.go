package docstore

import (
	"os/user"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/andrew-d/docstore/store"
)

var (
	mainCmd = &cobra.Command{
		Use:   "docstore",
		Short: "docstore is a way of saving, tagging, and organizing your documents",
	}

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new document store",
		Run:   runInit,
	}

	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List documents matching a certain query",
		Run:     wrapCommand(runList),
	}

	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show a single item (document, file, tag, etc.)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	flagVerbose       bool
	flagQuiet         bool
	flagStoreLocation string

	log = logrus.New()
)

func init() {
	usr, _ := user.Current()

	mainCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")
	mainCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "quiet output")
	mainCmd.PersistentFlags().StringVar(&flagStoreLocation, "path",
		filepath.Join(usr.HomeDir, ".docstore"), "location of the document store")
}

func Run() {
	// We defer adding top-level commands here so subcommands can register their
	// corresponding subcommands first (since init() functions execute in no
	// defined order).
	mainCmd.AddCommand(initCmd)
	mainCmd.AddCommand(listCmd)
	mainCmd.AddCommand(showCmd)

	mainCmd.Execute()
}

type CommandFunc func(cmd *cobra.Command, args []string, s *store.Store)

func wrapCommand(f CommandFunc) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if flagVerbose && flagQuiet {
			log.Fatal("Cannot be both verbose and quiet")
		}

		if flagVerbose {
			log.Level = logrus.DebugLevel
		} else if flagQuiet {
			log.Level = logrus.WarnLevel
		}

		st, err := store.Open(flagStoreLocation)
		if err != nil {
			log.WithField("err", err).Fatal("Could not open store")
		}

		f(cmd, args, st)
	}
}

func runList(cmd *cobra.Command, args []string, store *store.Store) {
	log.Info("Running list")
}

func runInit(cmd *cobra.Command, args []string) {
	_, err := store.New(flagStoreLocation)
	if err != nil {
		log.WithField("err", err).Error("Could not initialize docstore")
	} else {
		log.Info("Successfully initialized docstore")
	}
}
