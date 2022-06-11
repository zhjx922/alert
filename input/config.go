package input

import (
	"github.com/zhjx922/alert/output"
)

type Inputs struct {
	Name  string   `yaml:"name"`
	ScanFrequency int64 `yaml:"scan_frequency"`
	Paths []string `yaml:"paths"`
	IncludeLines []string `yaml:"include_lines"`
}

type Config struct {
	Inputs []*Inputs `yaml:"inputs"`
	OutputHttp *output.Http `yaml:"output.http"`
}
