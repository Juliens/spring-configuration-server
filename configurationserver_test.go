package springcongfigurationserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimpleConfiguration(t *testing.T) {
	var jsonContent string = `
	{"name":"application", "propertySources":[{"name":"test-name","source":{"test":"testvalue"}}]}

	`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, jsonContent)
	}))
	defer ts.Close()

	configurationServer := NewConfigurationServer(ts.URL)
	config := configurationServer.GetConfiguration("service", []string{"test"})

	testedValues := map[string]interface{
	}{
		"test":"testvalue",
	}

	for k, v := range testedValues {
		if configValue := config[k]; configValue != v {
			t.Logf("configVar %s must have %v but has: %v", k, v, configValue)
			t.Fail()
		}
	}

}


func TestMultipleProfileConfiguration(t *testing.T) {
	jsonContents := map[string]string {
		"/service-id/first":`
			{
			"name":"application",
			"propertySources":[
				{
					"name":"service-id",
					"source":{
					 	"var1":"first-service-id",
					 	"all-profiles":"first-service-id",
					 	"first-3-profiles":"first-service-id",
					 	"first-2-profiles":"first-service-id"
					}
				}, {
					"name":"application",
					"source":{
					 	"var2":"first-application",
					 	"all-profiles":"first-application",
					 	"first-3-profiles":"first-application",
					 	"first-2-profiles":"first-application",
					 	"first-profiles":"first-application"
					}
				}
			]}

			`,
		"/service-id/second":`
			{
			"name":"application",
			"propertySources":[
				{
					"name":"service-id",
					"source":{
					 	"var3":"second-service-id",
					 	"all-profiles":"second-service-id"
					}
				},
				{
					"name":"application",
					"source":{
					 	"var4":"second-application",
					 	"all-profiles":"second-application",
					 	"first-3-profiles":"second-application"
					}
				}
			]}
			`,
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if content, ok := jsonContents[r.URL.Path]; ok {
			fmt.Fprintln(w, content)
		} else {
			t.Logf("Wrong profile %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	configurationServer := NewConfigurationServer(ts.URL)
	config := configurationServer.GetConfiguration("service-id", []string{"first", "second"})

	existingConfigVars := []string{"var1", "var2", "var3", "var4"}

	for _, configVar := range existingConfigVars {
		if _, ok := config[configVar]; !ok {
			t.Logf("configVar %s must exist", configVar)
			t.Fail()
		}
	}

	testedValues := map[string]interface{
	}{
		"all-profiles":"second-service-id",
		"first-3-profiles":"second-application",
		"first-2-profiles":"first-service-id",
		"first-profiles":"first-application",
	}

	for k, v := range testedValues {
		if configValue := config[k]; configValue != v {
			t.Logf("configVar %s must have %v but has: %v", k, v, configValue)
			t.Fail()
		}
	}

}
