package service

import (
	"fmt"
	"github.com/abhinavk1/remote-config-server/pkg/util"
	"github.com/magiconair/properties"
	"log"
	"path/filepath"
)

type AbstractConfiguration interface {
	Get(applicationName, profile string) ([]byte, error)
}

type Configuration struct {
	workingDirectory string
}

func NewConfiguration(workingDirectory string) *Configuration {
	return &Configuration{
		workingDirectory: workingDirectory,
	}
}

func (svc *Configuration) Get(applicationName, profile string) ([]byte, error) {

	fileName := fmt.Sprintf("%s-%s.properties", applicationName, profile)
	filePath := filepath.Join(svc.workingDirectory, fileName)

	p, err := properties.LoadFile(filePath, properties.UTF8)
	if err != nil {
		log.Fatalf("error loading properties -> %v", err)
		return nil, err
	}

	jsonObject, err := util.PropertiesToJson(p.Map())
	if err != nil {
		log.Fatalf("error loading properties -> %v", err)
		return nil, err
	}

	return jsonObject, nil
}
