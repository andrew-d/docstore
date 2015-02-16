package main

import (
	"fmt"

	"github.com/mccoyst/validate"
)

var v validate.V

func init() {
	v = make(validate.V)
	v["nonzero"] = func(i interface{}) error {
		var ok bool

		switch v := i.(type) {
		case int:
			ok = v != 0
		case uint:
			ok = v != 0
		case string:
			ok = len(v) > 0
		}

		if !ok {
			return fmt.Errorf("should be nonzero")
		}

		return nil
	}
	v["nonempty"] = func(i interface{}) error {
		var ok bool

		switch v := i.(type) {
		case []int64:
			ok = len(v) > 0
		}

		if !ok {
			return fmt.Errorf("should be nonempty")
		}

		return nil
	}
	v["notnil"] = func(i interface{}) error {
		if i == nil {
			return fmt.Errorf("should not be nil")
		}
		return nil
	}
}
