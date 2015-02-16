package controllers

import (
	"fmt"
)

type VError struct {
	Base    interface{}
	Message string
	Status  int
}

func (e VError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Status, e.Message, e.Base)
}

type RedirectError struct {
	Location string
	Code     int
}

func (e RedirectError) Error() string {
	panic("Should not be used as an error")
}
