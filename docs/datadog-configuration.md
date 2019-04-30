# DataDog Synthetics Configuration

## Annotations

Additional datadog configurations can be added through a set of annotations to each ingress object, the current supported annotations are:

|                Annotation                |                      Description                       | Default |                   Valid Values                    |
| :--------------------------------------: | :----------------------------------------------------: | :-----: | :-----------------------------------------------: |
|    datadog.monitor.stakater.com/tags     | [Tags in Datadog](https://docs.datadoghq.com/tagging/) |  None   | Comma delimited string (e.g., `"tag1,tag2,tag3"`) |
| datadog.monitor.stakater.com/description |                  Monitor description                   |  None   |
|   datadog.monitor.stakater.com/paused    |               Pause or Resume a Monitor                |  false  |                 `true` or `false`                 |
