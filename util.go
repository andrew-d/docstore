package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zenazn/goji/web"
)

func errorsToStrings(errs []error) []string {
	ret := make([]string, 0, len(errs))
	for _, err := range errs {
		ret = append(ret, err.Error())
	}
	return ret
}

func parseIntParam(c web.C, name string) (int64, error) {
	val, found := c.URLParams[name]
	if !found {
		panic(fmt.Sprintf("no such parameter: '%s'", name))
	}

	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parameter '%s' is not numeric", name)
	}

	return num, nil
}

func parseUintParam(c web.C, name string) (uint64, error) {
	val, found := c.URLParams[name]
	if !found {
		panic(fmt.Sprintf("no such parameter: '%s'", name))
	}

	num, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parameter '%s' is not numeric", name)
	}

	return num, nil
}

func iQuery(s string) string {
	// TODO:
	if false {
		s = strings.TrimRight(s, "; ") + " RETURNING id"
	}

	return s
}
