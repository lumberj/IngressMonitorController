package datadog

const (
	// There are two possible monitor types, "api" and "browser". We currently only support "api".
	monitorType = "api"

	monitorTagsAnnotation   = "datadog.monitor.stakater.com/tags"
	monitorPausedAnnotation = "datadog.monitor.stakater.com/paused"

	// Required
	monitorDescriptionAnnotation = "datadog.monitor.stakater.com/description"

	// Possible values are
	// "aws:us-east-2", "aws:eu-central-1",
	// "aws:ca-central-1", "aws:eu-west-2", "aws:ap-northeast-1",
	// "aws:us-west-2", "aws:ap-southeast-2"
	monitorLocationsAnnotation = "datadog.monitor.stakater.com/locations"
	monitorLocationsDefault    = "aws:us-east-1"

	// Possible types are header, body, responseTime, and statusCode
	monitorResponseTypeAnnotation = "datadog.monitor.stakater.com/response-type"
	monitorResponseTypeDefault    = "statusCode"

	// Required if response-type is "header"
	monitorResponseHeaderKeyAnnotation   = "datadog.monitor.stakater.com/response-header-key"
	monitorResponseHeaderValueAnnotation = "datadog.monitor.stakater.com/response-header-value"
	// Valid values are "contains, does not contain, is, is not, matches, does not match"
	monitorResponseHeaderOperatorAnnotation = "datadog.monitor.stakater.com/response-header-operator"
	monitorResponseHeaderOperatorDefault    = "is"

	// Required if response-type is "responseTime"
	monitorResponseTimeLimitAnnotation = "datadog.monitor.stakater.com/response-time-limit"

	// Required if response-type is statusCode. Default is 200
	monitorResponseStatusCodeAnnotation = "datadog.monitor.stakater.com/response-status-code"
	monitorResponseStatusCodeDefault    = "200"

	// Required if response-type is "body"
	// Valid values are "contains, does not contain, is, is not, matches, does not match"
	monitorResponseBodyOperatorAnnotation = "datadog.monitor.stakater.com/response-body-operator"
	monitorResponseBodyOperatorDefault    = "is"
	monitorResponseBodyValueAnnotation    = "datadog.monitor.stakater.com/response-body"

	monitorRequestMethodAnnotation = "datadog.monitor.stakater.com/request-method"
	monitorRequestMethodDefault    = "GET"

	// Defaults found at https://docs.datadoghq.com/api/?lang=bash#create-a-test
	monitorRequestTimeoutAnnotation = "datadog.monitor.stakater.com/request-timeout"
	monitorRequestHeadersAnnotation = "datadog.monitor.stakater.com/request-headers"
	monitorRequestBodyAnnotation    = "datadog.monitor.stakater.com/request-body"

	// Current possible values are 60, 300, 900, 1800, 3600, 21600, 43200, 86400, 604800
	monitorRequestPeriodAnnotation = "datadog.monitor.stakater.com/request-period"
	monitorRequestPeriodDefault    = 60

	monitorRequestMinFailureDurationAnnotation = "datadog.monitor.stakater.com/request-min-failure-duration"
	monitorRequestMinFailureDurationDefault    = 0

	monitorMinLocationFailedAnnotation = "datadog.monitor.stakater.com/min-location-failed"
	monitorMinLocationFailedDefault    = 1

	monitorFollowRedirectsAnnotation = "datadog.monitor.stakater.com/follow-redirects"
	monitorFollowRedirectsDefault    = false
)
