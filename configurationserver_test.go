package springcongfigurationserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimpleConfiguration(t *testing.T) {
	var jsonContent string = `
	{"name":"application", "propertySources":[{"name":"test-name","source":{"server.port":8080}}]}

	`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, jsonContent)
	}))
	defer ts.Close()

	configurationServer := NewConfigurationServer(ts.URL)
	config := configurationServer.GetConfiguration([]string{"test"})
	if _, ok := config["server.port"]; !ok {
		t.Log("La configuration n'est pas récupérée")
		t.Fail()
	}
}


func TestMultipleProfileConfiguration(t *testing.T) {
	var firstJsonContent string = `
	{"name":"application", "propertySources":[{"name":"test-name","source":{"server.port":8080}}]}

	`
	var secondJsonContent string = `
	{"name":"application", "propertySources":[{"name":"test-name","source":{"server.port":8082, "server2.port":8080}}]}
	`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/application/first" {
			fmt.Fprintln(w, firstJsonContent)
		} else if r.URL.Path == "/application/second" {
			fmt.Fprintln(w, secondJsonContent)
		}
	}))
	defer ts.Close()

	configurationServer := NewConfigurationServer(ts.URL)
	config := configurationServer.GetConfiguration([]string{"first", "second"})

	if _, ok := config["server.port"]; !ok {
		t.Log("La configuration n'est pas récupérée")
		t.Fail()
	}
	if _, ok := config["server2.port"]; !ok {
		t.Log("La configuration n'est pas récupérée")
		t.Fail()
	}
	if (config["server.port"].(float64)) != 8082 {
		t.Logf("La valeur surchargée doit etre 8082 mais est %v", config["server.port"])
		t.Fail()
	}

	config = configurationServer.GetConfiguration([]string{"second", "first"})
	if _, ok := config["server.port"]; !ok {
		t.Log("La configuration n'est pas récupérée")
		t.Fail()
	}
	if _, ok := config["server2.port"]; !ok {
		t.Log("La configuration n'est pas récupérée")
		t.Fail()
	}
	if (config["server.port"].(float64)) != 8080 {
		t.Logf("La valeur surchargée doit etre 8080 mais est %v", config["server.port"])
		t.Fail()
	}
}
