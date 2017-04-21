package springcongfigurationserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type ConfigurationServer struct {
	endpoint string
}

func (configServer *ConfigurationServer) GetConfiguration(serviceId string, profiles []string) map[string]interface{} {
	toReturn := make(map[string]interface{})
	for _, profile := range profiles {
		toReturn = mergeMaps(toReturn, configServer.loadConfiguration(serviceId, profile))
	}
	return toReturn
}

func mergeMaps(original, added map[string]interface{}) map[string]interface{} {
	for k,v := range added {
		original[k] = v
	}
	return original
}

func (configServer *ConfigurationServer) loadConfiguration(serviceId, profile string) map[string]interface{} {
	url := fmt.Sprintf("%s/%s/%s", configServer.endpoint, serviceId, profile)
	log.Debugf("Reading configuration on %s", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Error: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error: %s", err)
	}
	configTest := configType{}
	err = json.Unmarshal(body, &configTest)
	if err != nil {
		log.Errorf("Error: %s", err)
	}
	toReturn := make(map[string]interface{})

	for i := len(configTest.PropertySources)-1 ; i >= 0; i-- {
		toReturn = mergeMaps(toReturn, configTest.PropertySources[i].Source)
	}
	return toReturn
}

func NewConfigurationServer(endpoint string) ConfigurationServer {
	return ConfigurationServer{endpoint: endpoint}
}

type propertySources struct {
	Source map[string]interface{} `json:"source"`
	Name   string                 `json:"name"`
}

type configType struct {
	Name            string            `json:"name"`
	PropertySources []propertySources `json:"propertySources"`
}
