package datadog

import (
	"log"
	"testing"

	"github.com/stakater/IngressMonitorController/pkg/config"
	"github.com/stakater/IngressMonitorController/pkg/models"
	"github.com/stakater/IngressMonitorController/pkg/util"
)

func makeService() *DatadogMonitorService {
	config := config.GetControllerConfig()
	service := DatadogMonitorService{}
	provider := util.GetProviderWithName(config, "Datadog")
	if nil == provider {
		panic("Failed to find provider")
	}
	service.Setup(*provider)
	return &service
}

func TestGetAllDatadogMonitors(t *testing.T) {
	service := makeService()
	monitors := service.GetAll()
	log.Printf("%+v\n", monitors)

	if len(monitors) == 0 {
		t.Log("No Monitors Exist")
	}
	if nil == monitors {
		t.Error("Error: " + "GetAll request Failed")
	}
}

func TestAddDatadogMonitorWithResponseType(t *testing.T) {

	annotations := make(map[string]string)
	// annotations["datadog.monitor.stakater.com/locations"] = "US-Central"

	m := models.Monitor{Name: "google-test", URL: "https://google.com", Annotations: annotations}
	service := makeService()
	service.Add(m)

	mRes, err := service.GetByName("google-test")

	if err != nil {
		t.Error("Error: " + err.Error())
	}
	if mRes.Name != m.Name {
		t.Error("The name is incorrect, expected: " + m.Name + ", but was: " + mRes.Name)
	}
	if mRes.URL != m.URL {
		t.Error("The URL is incorrect, expected: " + m.URL + ", but was: " + mRes.URL)
	}
	service.Remove(*mRes)
}
