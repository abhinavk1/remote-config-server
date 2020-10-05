package main

import (
	"github.com/abhinavk1/remote-config-server/pkg/configserver"
	"log"
)

func main() {

	s, err := configserver.New(&configserver.ServerConfig{
		PrivateKeyPath:   "",
		RepoUrl:          "",
		WorkingDirectory: "",
		Port:             8080,
	})

	if err != nil {
		log.Panic(err)
	}

	s.Start()
}
