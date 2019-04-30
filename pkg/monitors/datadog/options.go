package datadog

import (
	"strconv"

	"github.com/zorkian/go-datadog-api"
)

type wrappedDatadogSyntheticsTest struct {
	IsValid        bool
	Errors         []string
	SyntheticsTest *datadog.SyntheticsTest
}

func newWrappedDatadogSyntheticsTest(s *datadog.SyntheticsTest) *wrappedDatadogSyntheticsTest {
	w := wrappedDatadogSyntheticsTest{SyntheticsTest: s}
	w.IsValid = true
	return &w
}

func (w *wrappedDatadogSyntheticsTest) addError(err string) {
	w.Errors = append(w.Errors, err)
	w.IsValid = false
}

func (w *wrappedDatadogSyntheticsTest) setPaused(annotations map[string]string) {
}

func (w *wrappedDatadogSyntheticsTest) setDescription(annotations map[string]string) {
}

func (w *wrappedDatadogSyntheticsTest) setTags(annotations map[string]string) {}

func (w *wrappedDatadogSyntheticsTest) setLocations(annotations map[string]string) {
}

func (w *wrappedDatadogSyntheticsTest) setResponseType(annotations map[string]string) {
	responseType, ok := annotations[monitorResponseTypeAnnotation]
	if !ok {
		// Nothing to do. ResponseType not set.
		return
	}

	if responseType == "header" {
		headerKey, headerKeyExists := annotations[monitorResponseHeaderKeyAnnotation]
		headerValue, headerValueExists := annotations[monitorResponseHeaderValueAnnotation]

		operator, ok := annotations[monitorResponseHeaderOperatorAnnotation]
		if !ok {
			operator = monitorResponseHeaderOperatorDefault
		}

		if !headerKeyExists {
			err := ("Invalid response-type annotation. Missing Annotation: " + monitorResponseHeaderKeyAnnotation)
			w.addError(err)
			return
		}
		if !headerValueExists {
			err := ("Invalid response-type annotation. Missing Annotation: " + monitorResponseHeaderValueAnnotation)
			w.addError(err)
			return
		}
		w.SyntheticsTest.Config.Assertions = []datadog.SyntheticsAssertion{
			datadog.SyntheticsAssertion{
				Operator: &operator,
				Type:     &responseType,
				Property: &headerKey,
				Target:   &headerValue,
			},
		}
		return
	}

	if responseType == "body" {
		var responseBody string
		responseBody, ok := annotations[monitorResponseBodyValueAnnotation]
		if !ok {
			err := ("Invalid response-type annotation. Missing Annotation: " + monitorResponseBodyValueAnnotation)
			w.addError(err)
			return
		}

		operator, ok := annotations[monitorResponseBodyOperatorAnnotation]
		if !ok {
			operator = monitorResponseBodyOperatorDefault
		}
		w.SyntheticsTest.Config.Assertions = []datadog.SyntheticsAssertion{
			datadog.SyntheticsAssertion{
				Operator: &operator,
				Type:     &responseType,
				Target:   &responseBody,
			},
		}
		return
	}

	if responseType == "responseTime" {
		operator := "less than" // only possible value

		responseTime, ok := annotations[monitorResponseTimeLimitAnnotation]
		if !ok {
			err := ("Invalid response-type annotation. Missing Annotation: " + monitorResponseTimeLimitAnnotation)
			w.addError(err)
			return
		}

		responseTimeInt, err := strconv.Atoi(responseTime)
		if err != nil {
			err := ("Invalid response-time-limit annotation. Cannot cast to integer: " + responseTime)
			w.addError(err)
			return
		}

		w.SyntheticsTest.Config.Assertions = []datadog.SyntheticsAssertion{
			datadog.SyntheticsAssertion{
				Operator: &operator,
				Type:     &responseType,
				Target:   &responseTimeInt,
			},
		}
		return
	}

	if responseType == "statusCode" {
		operator := "is"
		statusCode, ok := annotations[monitorResponseStatusCodeAnnotation]
		if !ok {
			statusCode = monitorResponseStatusCodeDefault
		}

		statusCodeInt, err := strconv.Atoi(statusCode)
		if err != nil {
			err := ("Invalid response-status-code annotation. Cannot cast to integer: " + statusCode)
			w.addError(err)
			return
		}

		w.SyntheticsTest.Config.Assertions = []datadog.SyntheticsAssertion{
			datadog.SyntheticsAssertion{
				Operator: &operator,
				Type:     &responseType,
				Target:   &statusCodeInt,
			},
		}
	}

	err := ("Invalid response-type annotation: " + responseType)
	w.addError(err)
}

// func setResponseTimeLimit(s *datadog.SyntheticsTest, annotationValue string) {}

func (w *wrappedDatadogSyntheticsTest) SetOptionsFromAnnotations(annotations map[string]string) {
	annotationOptions := map[string]func(map[string]string){
		monitorPausedAnnotation:       w.setPaused,
		monitorDescriptionAnnotation:  w.setDescription,
		monitorTagsAnnotation:         w.setTags,
		monitorLocationsAnnotation:    w.setLocations,
		monitorResponseTypeAnnotation: w.setResponseType,
		// monitorRequestMethodAnnotation:             w.setRequestMethod,
		// monitorRequestTimeoutAnnotation:            w.setRequestTimeout,
		// monitorRequestHeadersAnnotation:            w.setRequestHeaders,
		// monitorRequestBodyAnnotation:               w.setRequestBody,
		// monitorRequestPeriodAnnotation:             w.setRequestPeriod,
		// monitorRequestMinFailureDurationAnnotation: w.setRequestMinFailureDuration,
		// monitorMinLocationFailedAnnotation:         w.setMinLocationFailed,
		// monitorFollowRedirectsAnnotation:           w.setFollowRedirects,
	}
	for annotation := range annotations {
		if fn, ok := annotationOptions[annotation]; ok {
			fn(annotations)
		}
	}
}

// var annotationOptions = map[string]func(s *datadog.SyntheticsTest, annotationValue string) error{
// 	monitorPausedAnnotation:       setPaused,
// 	monitorDescriptionAnnotation:  setDescription,
// 	monitorTagsAnnotation:         setTags,
// 	monitorLocationsAnnotation:    setLocations,
// 	monitorResponseTypeAnnotation: setResponseType,
// }

// func setOptions(s *datadog.SyntheticsTest, annotations map[string]string) {
// 	for annotation, value := range annotations {
// 		// validate()
// 		if fn, ok := annotationOptions[annotation]; ok {
// 			if err := fn(s, value); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// // TODO
// func setAnnotations(m *models.Monitor, s *datadog.SyntheticsTest) {}
