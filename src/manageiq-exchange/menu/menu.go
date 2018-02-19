package menu

import (
	"errors"
	"flag"
	"fmt"
	"manageiq-exchange/api"
	"manageiq-exchange/constants"
	"os"
)

type Configuration struct {
	Host      string
	Port      int
	Version   bool
	Providers bool
	Users     bool
	Expand    bool
}

func Menu() {
	config := &Configuration{}

	Banner()
	PassArguments(config)
	server, err := GetServer()
	if err != nil {
		fmt.Print(err)
	}

	var miqExchange api.Api
	miqExchange.Init(server, config.Port)

	statusConnection := miqExchange.CheckConnectionServer()
	if statusConnection {
		ShowInformationServer(config, miqExchange)
	}
}

func ShowInformationServer(configuration *Configuration, miqExchange api.Api) {
	if configuration.Version {
		info := miqExchange.GetInfo()
		fmt.Printf(info.Print())
	}

	if configuration.Users {
		users := miqExchange.GetUsers(configuration.Expand)
		fmt.Printf(users.Print(miqExchange.Data.Meta.TotalCount))
	}
}

func PassArguments(configuration *Configuration) {
	flag.StringVar(&configuration.Host, "host", "localhost", "specify host to use.  defaults to localhost.")
	flag.IntVar(&configuration.Port, "port", 0, "specify port to use.")
	flag.BoolVar(&configuration.Version, "version", false, "About version")
	flag.BoolVar(&configuration.Providers, "providers", false, "About providers")
	flag.BoolVar(&configuration.Users, "users", false, "About users")
	flag.BoolVar(&configuration.Expand, "expand", false, "Expand information")
	flag.Parse()
}

var myPrint = fmt.Println

func Banner() {
	myPrint("\033[0;31m", constants.BANNER, "\033[0m")
}

func GetOsEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func SetOsEnv(env string, key string) {
	os.Setenv(env, key)
}

func GetServer() (string, error) {
	server := GetOsEnv("EXCHANGE_SERVER", "")
	if len(server) == 0 {
		return "", errors.New("You need to set the environment EXCHANGE_SERVER (ex: localhost)")
	}
	return server, nil
}
