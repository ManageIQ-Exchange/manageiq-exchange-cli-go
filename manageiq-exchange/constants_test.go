package main

import (
	"reflect"
	"testing"
)

func TestConstants(t *testing.T) {
	var tests = []struct {
		message string
		input   string
		want    string
	}{
		{"SERVICE", SERVICE, "manageiq-exchange"},
		{"RESOURCE_PREFIX", RESOURCE_PREFIX, "//manageiq-exchange"},
		{"VERSION", VERSION, "dev"},
		{"BUILD_DATE", BUILD_DATE, "-"},
		{"SERVICEAPIVERSION", SERVICEAPIVERSION, "v0"},
	}
	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			if !reflect.DeepEqual(tt.input, tt.want) {
				t.Errorf("%s returned %+v, want %+v", tt.message, tt.input, tt.want)
			}
		})
	}
}
