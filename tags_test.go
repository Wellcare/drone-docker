package docker

import (
	"reflect"
	"testing"
)

func Test_stripTagPrefix(t *testing.T) {
	var tests = []struct {
		Before string
		After  string
	}{
		{"refs/tags/1.0.0", "1.0.0"},
		{"refs/tags/v1.0.0", "1.0.0"},
		{"v1.0.0", "1.0.0"},
	}

	for _, test := range tests {
		got, want := stripTagPrefix(test.Before), test.After
		if got != want {
			t.Errorf("Got tag %s, want %s", got, want)
		}
	}
}

func TestDefaultTags(t *testing.T) {
	var tests = []struct {
		Before string
		After  []string
	}{
		{"", []string{"latest"}},
		{"refs/heads/master", []string{"latest"}},
		{"refs/tags/0.9.0", []string{"0.9", "0.9.0"}},
		{"refs/tags/1.0.0", []string{"1", "1.0", "1.0.0"}},
		{"refs/tags/v1.0.0", []string{"1", "1.0", "1.0.0"}},
		{"refs/tags/v1.0.0-alpha.1", []string{"1.0.0-alpha.1"}},

		// malformed or errors
		{"refs/tags/x1.0.0", []string{"latest"}},
		{"v1.0.0", []string{"latest"}},
	}

	for _, test := range tests {
		got, want := DefaultTags(test.Before), test.After
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got tag %v, want %v", got, want)
		}
	}
}

func TestDefaultTagSuffix(t *testing.T) {
	var tests = []struct {
		Before string
		Suffix string
		After  []string
	}{
		// without suffix
		{
			After: []string{"latest"},
		},
		{
			Before: "refs/tags/v1.0.0",
			After: []string{
				"1",
				"1.0",
				"1.0.0",
			},
		},
		// with suffix
		{
			Suffix: "linux-amd64",
			After:  []string{"linux-amd64"},
		},
		{
			Before: "refs/tags/v1.0.0",
			Suffix: "linux-amd64",
			After: []string{
				"1-linux-amd64",
				"1.0-linux-amd64",
				"1.0.0-linux-amd64",
			},
		},
		{
			Suffix: "nanoserver",
			After:  []string{"nanoserver"},
		},
		{
			Before: "refs/tags/v1.9.2",
			Suffix: "nanoserver",
			After: []string{
				"1-nanoserver",
				"1.9-nanoserver",
				"1.9.2-nanoserver",
			},
		},
	}

	for _, test := range tests {
		got, want := DefaultTagSuffix(test.Before, test.Suffix), test.After
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got tag %v, want %v", got, want)
		}
	}
}
