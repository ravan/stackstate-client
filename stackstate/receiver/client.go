package receiver

import (
	"context"
	"crypto/tls"
	"fmt"
	rq "github.com/carlmjohnson/requests"
	sts "github.com/ravan/stackstate-client/stackstate"
	"golang.org/x/exp/maps"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	url      string
	conf     *sts.StackState
	instance *Instance
}

var (
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	DumpHttpRequest bool
)

const (
	Endpoint       string = "receiver/stsAgent/intake"
	MetricEndpoint string = "receiver/stsAgent/api/v1/series"
)

func NewClient(conf *sts.StackState, instance *Instance) *Client {
	url, _ := strings.CutSuffix(conf.ApiUrl, "/")
	return &Client{url: url, conf: conf, instance: instance}
}

func (c *Client) Send(f *Factory) error {
	if len(f.components) > 0 || len(f.events) > 0 {
		err := c.sendTopoAndEvents(f)
		if err != nil {
			return err
		}
	}

	if len(f.metrics) > 0 {
		series := MetricSeries{Series: f.metrics}
		if len(f.components) == 0 && len(f.events) == 0 {
			slog.Info("sending", "metrics", len(f.metrics))
		}
		return c.sendMetric(&series)
	}
	return nil
}

func (c *Client) sendMetric(series *MetricSeries) error {
	var e map[string]interface{}
	err := c.metricRequest().
		BodyJSON(series).
		ErrorJSON(&e).
		Fetch(context.Background())

	if err != nil {
		slog.Error("Failed to send metric data to receiver", slog.Any("error", err))
		return err
	}
	return nil
}

func (c *Client) sendTopoAndEvents(f *Factory) error {
	pl := NewEmptyStackStatePayload()
	pl.CollectionTimestamp = time.Now().Unix()
	pl.InternalHostname = f.source
	t := NewEmptyTopology()
	t.StartSnapshot = true
	t.StopSnapshot = true

	t.Components = maps.Values(f.components)
	t.Relations = maps.Values(f.relations)
	t.Instance.Type = c.instance.Type
	t.Instance.URL = c.instance.URL

	pl.Topologies = append(pl.Topologies, *t)
	if len(f.events) > 0 {
		pl.Events = map[string][]*Event{"events": f.events}
	} else {
		pl.Events = make(map[string][]*Event, 0)
	}

	slog.Info("sending", "components", len(t.Components),
		"relations", len(t.Relations), "events", len(f.events), "metrics", len(f.metrics))
	var e map[string]interface{}
	err := c.agentRequest().
		BodyJSON(&pl).
		ErrorJSON(&e).
		Fetch(context.Background())

	if err != nil {
		slog.Error("Failed to send data to receiver", "error", err, "details", e)
		return err
	}
	return nil
}

func (c *Client) agentRequest() *rq.Builder {
	uri := fmt.Sprintf("%s/%s", c.url, Endpoint)
	return request(uri).
		Param("api_key", c.conf.ApiKey)
}

func (c *Client) metricRequest() *rq.Builder {
	uri := fmt.Sprintf("%s/%s", c.url, MetricEndpoint)
	return request(uri).
		Param("api_key", c.conf.ApiKey)
}

func request(uri string) *rq.Builder {
	b := rq.URL(uri).
		ContentType("application/json")
	if DumpHttpRequest {
		b.Transport(rq.Record(nil, "http_dump"))
	} else {
		b.Transport(transport)
	}
	return b
}
