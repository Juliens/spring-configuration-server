package main

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

func (configServer *ConfigurationServer) GetConfiguration(profiles []string) map[string]interface{} {
	toReturn := make(map[string]interface{})
	for _, profile := range profiles {
		for k,v := range configServer.loadConfiguration(profile) {
			toReturn[k] = v
		}
	}
	return toReturn
}

func (configServer *ConfigurationServer) loadConfiguration(profile string) map[string]interface{} {
	url := fmt.Sprintf("%s/application/%s", configServer.endpoint, profile)
	log.Debugf("Lecture dans la config sur %s", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Impossible de lire la configuration dans le configserver, error: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Impossible de lire la configuration dans le configserver, error: %s", err)
	}
	configTest := configType{}
	json.Unmarshal(body, &configTest)
	if len(configTest.PropertySources)>0 {
		return configTest.PropertySources[0].Source
	}
	return make(map[string]interface{})
}

func copy (src map[string]interface{}, dest *interface{}) {
	*dest = src
}

func NewConfigurationServer(endpoint string) ConfigurationServer {
	return ConfigurationServer{endpoint: endpoint}
}

type Configuration struct {
	EurekaEndpoint string `json:"eureka.client.serviceUrl.defaultZone"`
	EurekaIPAddr   string `json:"eureka.instance.ipAddress"`
	EurekaPort     string `json:"eureka.instance.nonSecurePort"`
}

type propertySources struct {
	Source map[string]interface{} `json:"source"`
	Name   string                 `json:"name"`
}

type propertySourcesMap struct {
	Source map[string]string `json:"source"`
	Name   string            `json:"name"`
}

type configType struct {
	Name            string            `json:"name"`
	PropertySources []propertySources `json:"propertySources"`
}

type configMap struct {
	Name               string               `json:"name"`
	PropertySourcesMap []propertySourcesMap `json:"propertySources"`
}
