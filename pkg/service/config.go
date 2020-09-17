package service

type AbstractConfiguration interface {
	Get(applicationName, profile string) (interface{}, error)
}

type Configuration struct {
	workingDirectory string
}

func NewConfiguration(workingDirectory string) *Configuration {
	return &Configuration{
		workingDirectory: workingDirectory,
	}
}

func (svc *Configuration) Get(applicationName, profile string) (interface{}, error) {
	return nil, nil
}
