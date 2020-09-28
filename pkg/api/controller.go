package api

import (
	"errors"
	"github.com/abhinavk1/remote-config-server/pkg/service"
	"github.com/dimfeld/httptreemux/v5"
	"log"
	"net/http"
	"strings"
)

type Controller struct {
	configService service.AbstractConfiguration
}

func NewController(configService service.AbstractConfiguration) *Controller {
	return &Controller{
		configService: configService,
	}
}

func (controller *Controller) Handler(writer http.ResponseWriter, request *http.Request) {

	ctxData := httptreemux.ContextData(request.Context())
	params := ctxData.Params()
	param := params["param"]

	tag, extension := getTagAndExtension(param)

	switch extension {
	case "json":
		controller.jsonHandler(tag, writer)
		break

	case "properties":
		controller.propertiesHandler(tag, writer)
		break

	default:
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(""))
	}
}

func (controller *Controller) jsonHandler(tag string, writer http.ResponseWriter) {

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	applicationName, profile, err := getApplicationAndProfile(tag)
	if err != nil {
		log.Print(err)
		writer.Write([]byte("{}"))
		return
	}

	jsonConfig, err := controller.configService.GetJson(applicationName, profile)
	if err != nil {
		log.Print(err)
		writer.Write([]byte("{}"))
		return
	}

	writer.Write(jsonConfig)
}

func (controller *Controller) propertiesHandler(tag string, writer http.ResponseWriter) {

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)

	applicationName, profile, err := getApplicationAndProfile(tag)
	if err != nil {
		log.Print(err)
		writer.Write([]byte(""))
		return
	}

	properties, err := controller.configService.GetProperties(applicationName, profile)
	if err != nil {
		log.Print(err)
		writer.Write([]byte(""))
		return
	}

	writer.Write([]byte(properties))
}

func getApplicationAndProfile(tag string) (string, string, error) {
	var appName, profile string

	if tag == "" {
		return "", "", errors.New("invalid application name")
	}

	arr := strings.Split(tag, "-")
	appName = arr[0]

	if len(arr) > 1 {
		profile = arr[len(arr)-1]
	}

	return appName, profile, nil
}

func getTagAndExtension(param string) (string, string) {

	arr := strings.Split(param, ".")
	if len(arr) < 2 {
		return "", ""
	}

	return arr[0], strings.ToLower(arr[len(arr)-1])
}
