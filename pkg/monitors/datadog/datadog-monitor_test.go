package datadog

import (
	"log"
	"testing"

	"github.com/stakater/IngressMonitorController/pkg/config"
	"github.com/stakater/IngressMonitorController/pkg/util"
)

func TestGetAllDatadogMonitors(t *testing.T) {
	config := config.GetControllerConfig()

	service := DatadogMonitorService{}
	provider := util.GetProviderWithName(config, "Datadog")
	if provider == nil {
		panic("Failed to find provider")
	}
	service.Setup(*provider)
	monitors := service.GetAll()
	log.Println(monitors)

	if len(monitors) == 0 {
		t.Log("No Monitors Exist")
	}
	if nil == monitors {
		t.Error("Error: " + "GetAll request Failed")
	}
}
