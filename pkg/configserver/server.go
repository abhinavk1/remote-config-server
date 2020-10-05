package configserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/abhinavk1/remote-config-server/pkg/api"
	"github.com/abhinavk1/remote-config-server/pkg/router"
	"github.com/abhinavk1/remote-config-server/pkg/service"
	"github.com/abhinavk1/remote-config-server/pkg/versioncontrol"
	"github.com/go-http-utils/logger"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New(config *ServerConfig) (*Server, error) {

	if config == nil {
		return nil, errors.New("invalid config provided")
	}

	sshKey, err := ioutil.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	gitClient, err := versioncontrol.NewGitClient(&versioncontrol.GitClientConfig{
		Url:        config.RepoUrl,
		Path:       config.WorkingDirectory,
		PrivateKey: sshKey,
	})
	if err != nil {
		return nil, err
	}

	versionControlService := service.NewVersionControl(gitClient)

	go func() {
		pollError := versionControlService.PollRepository(20 * time.Second)
		if pollError != nil {
			log.Print(err)
		}
	}()

	//configWorkingDirectory := filepath.Join(config.WorkingDirectory, getRepoName(config.RepoUrl))
	configService := service.NewConfiguration(config.WorkingDirectory)
	controller := api.NewController(configService)

	app := router.New()
	app.GET("/:param", controller.Handler)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: logger.Handler(app, os.Stdout, logger.DevLoggerType),
	}

	return &Server{
		httpServer: httpServer,
	}, nil
}

func (server *Server) Start() error {
	return server.httpServer.ListenAndServe()
}

func (server *Server) Stop() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	// 5 seconds grace period before server shuts down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.httpServer.Shutdown(ctx); err != nil {
		log.Print(err)
	}
}

func getRepoName(repoUrl string) string {
	reg, _ := regexp.Compile(".*/")

	return reg.ReplaceAllString(repoUrl, "")
}
