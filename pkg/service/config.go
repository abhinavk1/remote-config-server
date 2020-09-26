package service

import (
	"fmt"
	"github.com/abhinavk1/remote-config-server/pkg/util"
	"github.com/magiconair/properties"
	"log"
	"path/filepath"
)

type AbstractConfiguration interface {
	GetJson(applicationName, profile string) ([]byte, error)
	GetProperties(applicationName, profile string) (string, error)
}

type Configuration struct {
	workingDirectory string
}

func NewConfiguration(workingDirectory string) *Configuration {
	return &Configuration{
		workingDirectory: workingDirectory,
	}
}

func (svc *Configuration) GetJson(applicationName, profile string) ([]byte, error) {

	profileProperties, err := svc.getPropertiesObject(applicationName, profile)
	if err != nil {
		log.Fatalf("error loading properties for application %v, profile %v -> %v",
			applicationName, profile, err)
		return nil, err
	}

	jsonObject, err := util.PropertiesToJson(profileProperties.Map())
	if err != nil {
		log.Fatalf("error converting properties to json for application %v, profile %v -> %v",
			applicationName, profile, err)
		return nil, err
	}

	return jsonObject, nil
}

func (svc *Configuration) GetProperties(applicationName, profile string) (string, error) {

	p, err := svc.getPropertiesObject(applicationName, profile)
	if err != nil {
		log.Fatalf("error loading properties for application %v, profile %v -> %v",
			applicationName, profile, err)
		return "", err
	}

	return p.String(), nil
}

func (svc *Configuration) getPropertiesObject(applicationName, profile string) (*properties.Properties, error) {

	profileConfigFilePath := svc.getFilePath(applicationName, profile)
	defaultConfigFilePath := svc.getFilePath(applicationName, "")

	filesToLoad := []string{profileConfigFilePath, defaultConfigFilePath}

	return properties.LoadFiles(filesToLoad, properties.UTF8, true)
}

func (svc *Configuration) getFilePath(applicationName, profile string) string {
	var fileName string

	if profile == "" {
		fileName = fmt.Sprintf("%s.properties", applicationName)
	} else {
		fileName = fmt.Sprintf("%s-%s.properties", applicationName, profile)
	}

	return filepath.Join(svc.workingDirectory, fileName)
}
