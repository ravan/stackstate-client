package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	rq "github.com/carlmjohnson/requests"
	sts "github.com/ravan/stackstate-client/stackstate"
	"log/slog"
	"net/http"
	"strings"
)

type Client struct {
	url  string
	conf *sts.StackState
}

var (
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	DumpHttpRequest bool
)

const (
	GroovyScript string = "GroovyScript"
)

func NewClient(conf *sts.StackState) *Client {
	url, _ := strings.CutSuffix(conf.ApiUrl, "/")
	return &Client{url: url, conf: conf}
}

func (c *Client) Status() (*ServerInfo, error) {
	var s ServerInfo
	err := c.apiRequests("server/info").
		ToJSON(&s).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *Client) ViewSnapshot(req *ViewSnapshotRequest) (*ViewSnapshotResponse, error) {
	var res ViewSnapshotResponse
	err := c.apiRequests("snapshot").
		Post().
		BodyJSON(&req).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) Layers() (*map[int64]NodeType, error) {
	return c.getNodesOfType("Layer")
}

func (c *Client) ComponentTypes() (*map[int64]NodeType, error) {
	return c.getNodesOfType("ComponentType")
}

func (c *Client) RelationTypes() (*map[int64]NodeType, error) {
	return c.getNodesOfType("RelationType")
}

func (c *Client) Domains() (*map[int64]NodeType, error) {
	return c.getNodesOfType("Domain")
}

func (c *Client) getNodesOfType(t string) (*map[int64]NodeType, error) {
	var res []NodeType
	err := c.apiRequests(fmt.Sprintf("node/%s", t)).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	nodes := make(map[int64]NodeType, len(res))
	for _, r := range res {
		nodes[r.ID] = r
	}
	return &nodes, nil
}

func (c *Client) TopologyQuery(query string, at string, fullLoad bool) (*TopoQueryResponse, error) {
	query, at = sanitizeQuery(query, at)
	method := "components"
	if fullLoad {
		method = "fullComponents"
	}
	body := fmt.Sprintf(`Topology.query('%s')%s.%s()`, query, at, method)
	return c.executeTopoScript(scriptRequest{
		ReqType: GroovyScript,
		Body:    body,
	})
}

func (c *Client) TopologyStreamQuery(query string, at string, withSyncData bool) (*TopoQueryResponse, error) {
	query, at = sanitizeQuery(query, at)
	method := ""
	if withSyncData {
		method = ".withSynchronizationData()"
	}
	body := fmt.Sprintf(`TopologyStream.query('%s')%s%s`, query, at, method)
	return c.executeTopoScript(scriptRequest{
		ReqType: GroovyScript,
		Body:    body,
	})
}

func sanitizeQuery(query string, at string) (string, string) {
	query = strings.ReplaceAll(query, "'", "\"")
	if at != "" {
		at = fmt.Sprintf(".at('%s')", at)
	}
	return query, at
}

func (c *Client) executeTopoScript(req scriptRequest) (*TopoQueryResponse, error) {
	var r SuccessResp
	var e ErrorResp
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	slog.Debug("request", "body", string(b))
	err = c.apiRequests("script").
		BodyJSON(&req).
		ErrorJSON(&e).
		ToJSON(&r).
		Fetch(context.Background())
	if err != nil {
		if e.Errors != nil {
			return &TopoQueryResponse{Success: false, Errors: e.Errors, Data: nil}, nil
		}
		return nil, err
	}
	return &TopoQueryResponse{Success: true, Errors: nil, Data: r.Result}, nil
}

func (c *Client) apiRequests(endpoint string) *rq.Builder {
	uri := fmt.Sprintf("%s/api/%s", c.url, endpoint)
	return request(uri).
		Header("X-API-Token", c.conf.ApiToken)
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
