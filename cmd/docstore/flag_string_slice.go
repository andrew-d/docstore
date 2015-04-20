package docstore

import (
	"strings"

	"github.com/spf13/pflag"
)

type StringSlice []string

func (s *StringSlice) String() string {
	return "[" + strings.Join(*s, ",") + "]"
}

func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s *StringSlice) Type() string {
	return "StringSlice"
}

var _ pflag.Value = &StringSlice{}
