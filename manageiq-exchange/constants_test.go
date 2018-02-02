package main

import (
	"reflect"
	"testing"
)

func TestConstants(t *testing.T) {

	want := "manageiq-exchange"
	if !reflect.DeepEqual(SERVICE, want) {
		t.Errorf("SERVICE returned %+v, want %+v", SERVICE, want)
	}

	want = "//manageiq-exchange"
	if !reflect.DeepEqual(RESOURCE_PREFIX, want) {
		t.Errorf("RESOURCE_PREFIX returned %+v, want %+v", RESOURCE_PREFIX, want)
	}

	want = "dev"
	if !reflect.DeepEqual(VERSION, want) {
		t.Errorf("VERSION returned %+v, want %+v", VERSION, want)
	}

	want = "-"
	if !reflect.DeepEqual(BUILD_DATE, want) {
		t.Errorf("BUILD_DATE returned %+v, want %+v", BUILD_DATE, want)
	}

	want = "v0"
	if !reflect.DeepEqual(SERVICEAPIVERSION, want) {
		t.Errorf("SERVICEAPIVERSION returned %+v, want %+v", SERVICEAPIVERSION, want)
	}
}
