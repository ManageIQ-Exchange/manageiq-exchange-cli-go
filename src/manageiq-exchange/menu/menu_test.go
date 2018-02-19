package menu

import (
	"manageiq-exchange/api"
	"manageiq-exchange/constants"
	"os"
	"reflect"
	"testing"
)

func TestGetOsEnvAndSetosEnv(t *testing.T) {
	want := "manageiq-exchange"
	os.Setenv("TEST_EXCHANGE_ENV", "manageiq-exchange")
	if !reflect.DeepEqual(GetOsEnv("TEST_EXCHANGE_ENV", "another_env"), want) {
		t.Errorf("GetOsEnv returned %+v, want %+v", GetOsEnv("TEST_EXCHANGE_ENV", "another_env"), want)
	}
	want = "another_env"
	if !reflect.DeepEqual(GetOsEnv("TEST_EXCHANGE_ENV_SECOND", "another_env"), want) {
		t.Errorf("GetOsEnv returned %+v, want %+v", GetOsEnv("TEST_EXCHANGE_ENV_SECOND", "another_env"), want)
	}

	want = "GO_ROCKS"
	SetOsEnv("TEST_EXCHANGE_ENV", "GO_ROCKS")
	if !reflect.DeepEqual(os.Getenv("TEST_EXCHANGE_ENV"), want) {
		t.Errorf("SetOsEnv returned %+v, want %+v", os.Getenv("TEST_EXCHANGE_ENV"), want)
	}
}

func TestGetServer(t *testing.T) {
	os.Unsetenv("EXCHANGE_SERVER")
	want := ""
	serv, err := GetServer()
	if !reflect.DeepEqual(serv, want) {
		t.Errorf("GetServer returned %+v, want %+v", serv, want)
	}
	want = "You need to set the environment EXCHANGE_SERVER (ex: localhost)"
	if !reflect.DeepEqual(err.Error(), want) {
		t.Errorf("GetServer returned %+v, want %+v", err.Error(), want)
	}

	os.Setenv("EXCHANGE_SERVER", "localhost")

	want = "localhost"
	serv, err = GetServer()
	if !reflect.DeepEqual(serv, want) {
		t.Errorf("GetServer returned %+v, want %+v", serv, want)
	}
}

func TestPassArguments(t *testing.T) {
	gotConfig := &Configuration{}
	wantConfig := &Configuration{
		Host:      "localhost",
		Port:      0,
		Version:   false,
		Providers: false,
		Users:     false,
		Expand:    false,
	}
	PassArguments(gotConfig)
	if !reflect.DeepEqual(gotConfig, wantConfig) {
		t.Errorf("ReadFlags(gotConfig) returned %+v, want %+v", gotConfig, wantConfig)
	}
}

func TestShowInformationServer(t *testing.T) {
	config := &Configuration{}
	miqExchange := api.Api{}
	ShowInformationServer(config, miqExchange)
}

func TestBanner(t *testing.T) {
	oldMyPrint := myPrint
	defer func() { myPrint = oldMyPrint }()

	var gotOutput string
	fakePrint := func(s ...interface{}) (n int, err error) {
		for i := range s {
			gotOutput += s[i].(string)
		}
		return len(s), nil
	}

	myPrint = fakePrint
	wantBanner := "\033[0;31m" + constants.BANNER + "\033[0m"
	Banner()

	if !reflect.DeepEqual(gotOutput, wantBanner) {
		t.Errorf("Banner() returned %+v, want %+v", gotOutput, wantBanner)
	}
}
