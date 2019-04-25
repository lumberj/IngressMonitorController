package datadog

import (
	"fmt"
	"github.com/stakater/IngressMonitorController/pkg/config"
	"github.com/stakater/IngressMonitorController/pkg/models"
	"github.com/zorkian/go-datadog-api"
	"log"
)

// MonitorService for Datadog Synthetics
type MonitorService struct {
	apiKey string
	appKey string
	url    string
	client *datadog.Client
}

// Setup DataDog Synthetics MonitorService
func (service *MonitorService) Setup(p config.Provider) {
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
	// 	log.Println("Error initializing monitor: ", err.Error())
	// 	return
	// }

}

// GetAll retrieve all Datadog Synthetic Test Monitors
func (service *MonitorService) GetAll() []models.Monitor {
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

// Add new Monitor to Datadog Synthetics Tests
func (service *MonitorService) Add(m models.Monitor) {
	var (
		method       = "GET"
		responseType = "statusCode"
		target       = 200
		operator     = "is"
		locations    = []string{"aws:us-east-1"}
		tickEvery    = 60 // seconds
		testType     = "api"
	)
	synRequest := datadog.SyntheticsRequest{
		Url:    &m.URL,
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
		Config:    &config,
		Locations: locations,
		Options:   &options,
		Type:      &testType,
	}

	if _, err := service.client.CreateSyntheticsTest(&synTest); err != nil {
		log.Println("Error received creating Synthetics Tests: ", err.Error())
		return
	}
	log.Println("Monitor Added: " + m.Name)
}

// Update Datadog Synthetics Tests Monitor
func (service *MonitorService) Update(m models.Monitor) {
	synTest, err := service.client.GetSyntheticsTest(m.ID)
	if err != nil {
		log.Println("Error getting monitor: "+m.Name, err.Error())
	}
	// TODO: Allow updating other options based on Annotations
	synTest.Name = &m.Name
	synTest.Config.Request.Url = &m.URL
}

// GetByName returns Datadog Synthetics Test Monitor by Name
func (service *MonitorService) GetByName(n string) (*models.Monitor, error) {
	var monitor *models.Monitor
	syntheticsTests, err := service.client.GetSyntheticsTests()
	if err != nil {
		log.Println("Error received while listing Synthetics Tests: ", err.Error())
		return monitor, nil
	}
	for _, synTest := range syntheticsTests {
		if *synTest.Name == n {
			monitor = &models.Monitor{
				ID:   *synTest.PublicId,
				Name: *synTest.Name,
				URL:  *synTest.Config.Request.Url,
			}
			return monitor, nil
		}
	}
	if monitor == nil {
		return nil, fmt.Errorf("Could not find monitor with name: %v", n)
	}
	return nil, nil
}

// Remove a Datadog Synthetics Test Monitor
func (service *MonitorService) Remove(m models.Monitor) {
	err := service.client.DeleteSyntheticsTests([]string{m.ID})
	if err != nil {
		log.Println("Error deleting monitor: ", m.Name, err.Error())
	}
}