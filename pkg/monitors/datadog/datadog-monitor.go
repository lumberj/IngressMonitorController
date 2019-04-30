package datadog

import (
	"fmt"
	"log"

	"github.com/stakater/IngressMonitorController/pkg/config"
	"github.com/stakater/IngressMonitorController/pkg/models"
	"github.com/zorkian/go-datadog-api"
)

// DatadogMonitorService for Datadog Synthetics
type DatadogMonitorService struct {
	apiKey string
	appKey string
	url    string
	client *datadog.Client
}

func DatadogMonitorToBaseMonitorMapper(datadogMonitor *datadog.SyntheticsTest) *models.Monitor {
	monitor := &models.Monitor{
		ID:   *datadogMonitor.PublicId,
		Name: *datadogMonitor.Name,
		URL:  *datadogMonitor.Config.Request.Url,
	}
	// TODO
	// setAnnotations(monitor, datadogMonitor)
	return monitor
}

func DatadogMonitorToBaseMonitorMappers(datadogMonitors *[]datadog.SyntheticsTest) *[]models.Monitor {
	monitors := make([]models.Monitor, len(*datadogMonitors))
	for i, synTest := range *datadogMonitors {
		newMon := DatadogMonitorToBaseMonitorMapper(&synTest)
		monitors[i] = *newMon
	}
	return &monitors
}

// Setup DataDog Synthetics MonitorService
func (service *DatadogMonitorService) Setup(p config.Provider) {
	service.apiKey = p.ApiKey
	service.appKey = p.AppKey
	service.url = p.ApiURL

	service.client = datadog.NewClient(service.apiKey, service.appKey)

	if len(service.url) > 0 {
		service.client.SetBaseUrl(service.url)
	}

	// TODO: Validate Connection
	// _, err := service.client.Validate()
	// if err != nil {
	// 	log.Println("Error initializing Datadog monitor: ", err.Error())
	// 	return
	// }

}

// GetAll retrieve all Datadog Synthetic Test Monitors
func (service *DatadogMonitorService) GetAll() []models.Monitor {
	syntheticsTests, err := service.client.GetSyntheticsTests()

	if err != nil {
		log.Println("Error received while listing Synthetics Tests: ", err.Error())
		return nil
	}

	monitors := make([]models.Monitor, len(syntheticsTests))

	for i, synTest := range syntheticsTests {
		newMon := models.Monitor{
			ID:   *synTest.PublicId,
			Name: *synTest.Name,
			URL:  *synTest.Config.Request.Url,
		}
		monitors[i] = newMon
	}
	return monitors
}

func monitorToSyntheticsTest(m models.Monitor) (*datadog.SyntheticsTest, error) {
	// Defaults
	var (
		url          = m.URL
		name         = m.Name
		message      = m.Name
		method       = monitorRequestMethodDefault
		responseType = monitorResponseTypeDefault
		target       = monitorResponseStatusCodeDefault
		operator     = "is"
		locations    = []string{monitorLocationsDefault}
		tickEvery    = monitorRequestPeriodDefault
		testType     = monitorType
	)

	synRequest := datadog.SyntheticsRequest{
		Url:    &url,
		Method: &method,
	}

	assertions := []datadog.SyntheticsAssertion{
		datadog.SyntheticsAssertion{
			Operator: &operator,
			Type:     &responseType,
			Target:   &target,
		},
	}

	config := datadog.SyntheticsConfig{
		Request:    &synRequest,
		Assertions: assertions,
	}

	options := datadog.SyntheticsOptions{
		TickEvery: &tickEvery,
	}

	synTest := datadog.SyntheticsTest{
		Name:      &name,
		Message:   &message,
		Config:    &config,
		Locations: locations,
		Options:   &options,
		Type:      &testType,
	}
	wrapped := newWrappedDatadogSyntheticsTest(&synTest)
	wrapped.SetOptionsFromAnnotations(m.Annotations)
	if !wrapped.IsValid {
		return nil, fmt.Errorf("Error creating monitor: %v"+m.Name, wrapped.Errors)
	}

	return &synTest, nil
}

// Add new Monitor to Datadog Synthetics Tests
func (service *DatadogMonitorService) Add(m models.Monitor) {
	synTest, err := monitorToSyntheticsTest(m)
	if err != nil {
		log.Println("Error creating monitor: "+m.Name, err.Error())
		return
	}

	if _, err := service.client.CreateSyntheticsTest(synTest); err != nil {
		log.Println("Error received creating Synthetics Test for Monitor: "+m.Name, err.Error())
		return
	}
	log.Println("Monitor Added: " + m.Name)
}

// Update Datadog Synthetics Tests Monitor
func (service *DatadogMonitorService) Update(m models.Monitor) {
	newSyntest, _ := monitorToSyntheticsTest(m)
	if _, err := service.client.UpdateSyntheticsTest(m.ID, newSyntest); err != nil {
		log.Println("Error updating monitor: "+m.Name, err.Error())
	}
}

// GetByName returns Datadog Synthetics Test Monitor by Name
func (service *DatadogMonitorService) GetByName(n string) (*models.Monitor, error) {
	var monitor *models.Monitor
	syntheticsTests, err := service.client.GetSyntheticsTests()
	if err != nil {
		log.Println("Error received while listing Synthetics Tests: ", err.Error())
		return monitor, nil
	}
	for _, synTest := range syntheticsTests {
		if *synTest.Name == n {
			monitor = DatadogMonitorToBaseMonitorMapper(&synTest)
			return monitor, nil
		}
	}
	if monitor == nil {
		return nil, fmt.Errorf("Could not find monitor with name: %v", n)
	}
	return nil, nil
}

// Remove a Datadog Synthetics Test Monitor
func (service *DatadogMonitorService) Remove(m models.Monitor) {
	err := service.client.DeleteSyntheticsTests([]string{m.ID})
	if err != nil {
		log.Println("Error deleting monitor: ", m.Name, err.Error())
	}
}
