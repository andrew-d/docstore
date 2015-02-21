package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji/web"
)

// Helper type for a JSON-style map
type M map[string]interface{}

// Helper function to render JSON to a http.ResponseWriter
func (c *AppController) JSON(w http.ResponseWriter, status int, val interface{}) error {
	var result []byte
	var err error

	if c.Debug {
		result, err = json.MarshalIndent(val, "", "  ")
		result = append(result, '\n')
	} else {
		result, err = json.Marshal(val)
	}

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, err = w.Write(result)
	return err
}

// Helper function to parse an integer parameter
func (c *AppController) parseIntParam(ctx web.C, name string) (int64, error) {
	val, found := ctx.URLParams[name]
	if !found {
		panic(fmt.Sprintf("no such parameter: '%s'", name))
	}

	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, VError{
			Base:    fmt.Errorf("parameter '%s' is not numeric", name),
			Message: fmt.Sprintf("invalid parameter: %s", name),
			Status:  400,
		}
	}

	return num, nil
}

// Helper function to (possibly) add `RETURNING id` to a query, if we're
// using a Postgres database.
func (c *AppController) iQuery(s string) string {
	if c.DB.DriverName() == "postgres" {
		s = strings.TrimRight(s, "; ") + " RETURNING id"
	}

	return s
}

// Helper function to run some code within a transaction, properly committing
// or rolling back depending on the return value.
func (c *AppController) inTransaction(cb func(tx *sqlx.Tx) error) error {
	finished := false
	tx, err := c.DB.Beginx()
	if err != nil {
		return VError{err, "error creating transaction", http.StatusInternalServerError}
	}

	defer func() {
		if !finished {
			tx.Rollback()
		}
	}()

	err = cb(tx)
	if err != nil {
		return err
	}

	finished = true
	err = tx.Commit()
	if err != nil {
		return VError{err, "error committing transaction", http.StatusInternalServerError}
	}

	return nil
}
