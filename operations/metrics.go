package operations

import (
	"fmt"
	"net/url"
	"time"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/xviper"
)

const (
	trackingUrl = `/metric-v1/%v/%v/%v/%v/%v`
	metricsHost = `https://telemetry.robocorp.com`
)

func sendMetric(kind, name, value string) {
	client, err := cloud.NewClient(metricsHost)
	if err != nil {
		common.Debug("ERROR: %v", err)
		return
	}
	timestamp := time.Now().UnixNano()
	url := fmt.Sprintf(trackingUrl, url.PathEscape(kind), timestamp, url.PathEscape(xviper.TrackingIdentity()), url.PathEscape(name), url.PathEscape(value))
	common.Debug("DEBUG: Sending metric as %v%v", metricsHost, url)
	client.Put(client.NewRequest(url))
}

func SendMetric(kind, name, value string) {
	common.Debug("DEBUG: SendMetric kind:%v name:%v value:%v send:%v", kind, name, value, xviper.CanTrack())
	if xviper.CanTrack() {
		sendMetric(kind, name, value)
	}
}

func BackgroundMetric(kind, name, value string) {
	go SendMetric(kind, name, value)
}
