package api

import (
	"fmt"
	"manageiq-exchange/models/info"
	meta "manageiq-exchange/models/metadata"
	user "manageiq-exchange/models/user"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the client being tested.
	client *API

	// urlTest is the url test server
	urlTest *url.URL

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a httpClient that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// http client configured to use test server
	urlTest, _ = url.Parse(server.URL + "/")
	i, _ := strconv.Atoi(urlTest.Port())
	client = &API{}
	client.Init(urlTest.Hostname(), i)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestApi_Init(t *testing.T) {
	var server API
	inputServer := "localhost"
	inputPort := 3000
	wantServer := "localhost"
	wantPort := 3000
	server.Init(inputServer, inputPort)
	if !reflect.DeepEqual(server.Server, wantServer) || !reflect.DeepEqual(server.Port, wantPort) {
		t.Errorf("Api.Init(%v, %v) returned server=%v and port=%v, want server=%v and port=%v", inputServer, inputPort, server.Server, server.Port, wantServer, wantPort)
	}
}

func TestApi_CheckConnectionServer(t *testing.T) {
	want := true
	setup()
	statusConn := client.CheckConnectionServer()
	teardown()
	if statusConn != want {
		t.Errorf("Api.CheckConnectionServer() returned %v, want %v", statusConn, want)
	}
}

func TestApi_CheckConnectionServer_KO(t *testing.T) {
	setup()
	want := false
	client.Port = 0
	statusConn := client.CheckConnectionServer()
	teardown()
	if statusConn != want {
		t.Errorf("Api.CheckConnectionServer() returned %v, want %v", statusConn, want)
	}
}

func TestApi_URL(t *testing.T) {
	var tests = []struct {
		inputServer string
		inputPort   int
		wantURL     string
	}{
		{"localhost", 0, "http://localhost"},
		{"localhost", 3000, "http://localhost:3000"},
	}
	for _, tt := range tests {
		t.Run(tt.inputServer, func(t *testing.T) {
			var server API
			server.Init(tt.inputServer, tt.inputPort)
			gotURL := server.URL()
			if !reflect.DeepEqual(gotURL, tt.wantURL) {
				t.Fatalf("Api.URL() returned %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func TestApi_Request(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"Data":{"A":"a"}}`)
	})

	err := client.Request(http.MethodGet, "/", nil)

	if err != nil {
		t.Errorf("Not expected errors to be returned")
	}

	wantData := map[string]interface{}{
		"A": "a",
	}
	wantMeta := meta.Metadata{
		CurrentPage: 0,
		TotalPages:  0,
		TotalCount:  0,
	}
	if !reflect.DeepEqual(client.Data.Data, wantData) {
		t.Errorf("client.Data.Data returned %+v want %+v", client.Data.Data, wantData)
	}

	if !reflect.DeepEqual(client.Data.Meta, wantMeta) {
		t.Errorf("client.Data.Meta returned %+v want %+v", client.Data.Meta, wantMeta)
	}
}

func TestApi_Request_badURL(t *testing.T) {
	setup()
	defer teardown()

	err := client.Request(http.MethodGet, "%zzzzz", nil)

	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestApi_Request_invalidDo(t *testing.T) {
	setup()
	defer teardown()

	client.Port = 0

	err := client.Request(http.MethodGet, ":", nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestApi_Request_httpError(t *testing.T) {
	setup()
	defer teardown()

	pathURL := "/foo"

	mux.HandleFunc(pathURL, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	err := client.Request(http.MethodPost, pathURL, nil)

	if err == nil {
		t.Error("Expected HTTP 400 errors.")
	}
}

func TestApi_GetInfo(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"data":{"version":"1.0","providers":{"github.com":{"server":"github.com","version":"v3","id_application":"abc"}}}}`)
	})

	gotInfo := client.GetInfo()
	wantInfo := info.Info{
		Version: "1.0",
		Providers: map[string]info.Provider{
			"github.com": info.Provider{
				ApplicationID: "abc",
				Server:        "github.com",
				Version:       "v3",
			},
		},
	}

	if !reflect.DeepEqual(gotInfo, wantInfo) {
		t.Errorf("Api.GetInfo() returned %+v want %+v", gotInfo, wantInfo)
	}
}

func TestApi_GetInfo_httpError(t *testing.T) {
	setup()
	defer teardown()

	pathURL := "/"

	mux.HandleFunc(pathURL, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	gotInfo := client.GetInfo()
	wantInfo := info.Info{}

	if !reflect.DeepEqual(gotInfo, wantInfo) {
		t.Errorf("Api.GetInfo() returned %+v want %+v", gotInfo, wantInfo)
	}
}

func TestApi_GetUsers_withoutExpandResources(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"data":[{"login":"aljesusg","name":"Alberto","github_id":1}]}`)
	})

	expand := false

	gotUsers := client.GetUsers(expand)
	wantUser := user.User{
		Login:    "aljesusg",
		Name:     "Alberto",
		GithubID: 1,
	}
	wantUsers := user.UserCollection{}
	wantUsers.Users = append(wantUsers.Users, wantUser)
	wantUsers.Total = len(wantUsers.Users)

	if !reflect.DeepEqual(gotUsers, wantUsers) {
		t.Errorf("Api.GetUsers(%t) returned %+v want %+v", expand, gotUsers, wantUsers)
	}
}

func TestApi_GetUsers_withExpandResources(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"data":[{"login":"aljesusg","name":"Alberto","github_id":1}]}`)
	})

	expand := true

	gotUsers := client.GetUsers(expand)
	wantUser := user.User{
		Login:    "aljesusg",
		Name:     "Alberto",
		GithubID: 1,
	}
	wantUsers := user.UserCollection{}
	wantUsers.Users = append(wantUsers.Users, wantUser)
	wantUsers.Total = len(wantUsers.Users)

	if !reflect.DeepEqual(gotUsers, wantUsers) {
		t.Errorf("Api.GetUsers(%t) returned %+v want %+v", expand, gotUsers, wantUsers)
	}
}

func TestApi_GetUsers_httpError(t *testing.T) {
	setup()
	defer teardown()

	pathURL := "/v1/users"
	expand := true

	mux.HandleFunc(pathURL, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	gotUsers := client.GetUsers(expand)
	wantUsers := user.UserCollection{}

	if !reflect.DeepEqual(gotUsers, wantUsers) {
		t.Errorf("Api.GetUsers(%t) returned %+v want %+v", expand, gotUsers, wantUsers)
	}
}
