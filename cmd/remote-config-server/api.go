package main

import (
	"github.com/abhinavk1/remote-config-server/pkg/api"
	"github.com/abhinavk1/remote-config-server/pkg/router"
	"github.com/abhinavk1/remote-config-server/pkg/service"
	"github.com/abhinavk1/remote-config-server/pkg/versioncontrol"
	"github.com/dimfeld/httptreemux/v5"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func newServer() *httptreemux.ContextMux {

	privateKeyPath := os.Getenv("REPO_PRIVATE_KEY_PATH")
	sshKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Panic(err)
	}

	gitClient, err := versioncontrol.NewGitClient(&versioncontrol.GitClientConfig{
		Url:        os.Getenv("REPO_URL"),
		Path:       os.Getenv("REPO_CLONE_PATH"),
		PrivateKey: sshKey,
	})
	if err != nil {
		log.Panic(err)
	}

	versionControlService := service.NewVersionControl(gitClient)

	go func() {
		pollError := versionControlService.PollRepository(20 * time.Second)
		if pollError != nil {
			log.Print(err)
		}
	}()

	configService := service.NewConfiguration(os.Getenv("WORKING_DIRECTORY"))
	controller := api.NewController(configService)

	app := router.New()
	app.GET("/:param", controller.Handler)

	return app
}
