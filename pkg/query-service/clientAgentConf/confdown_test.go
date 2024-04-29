package clientAgentConf

import (
	"fmt"
	"testing"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/collector/receiver"
	"gopkg.in/yaml.v3"
)

func Test(t *testing.T) {
	data := map[string]any{}
	err := yaml.Unmarshal([]byte(`receivers:
	filelog/lastlog:
	  include: [ /var/log/lastlog ]
	  start_at: end
  service:
	pipelines:
	  logs:
		receivers: [filelog/lastlog]
		processors: [batch, attributes]
		exporters: [ otlp/log ]`), &data)
	if err != nil {
		fmt.Println(err)
	}
	receiver.
	service.New()
	service.Config
	c := confmap.NewFromStringMap(data)

}
