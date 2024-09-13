package api

import (
	"encoding/json"
	"time"
)

// StackState Server Api DTOs

type ServerInfo struct {
	Version struct {
		Major  int    `json:"major"`
		Patch  int    `json:"patch"`
		Diff   string `json:"diff"`
		Commit string `json:"commit"`
		IsDev  bool   `json:"isDev"`
	} `json:"version"`
	DeploymentMode string `json:"deploymentMode"`
}

type SyncComponent struct {
	Id                  int                    `json:"id"`
	Identifiers         []string               `json:"identifiers"`
	Labels              []Label                `json:"labels"`
	Environments        []int                  `json:"environments"`
	Domain              int                    `json:"domain"`
	LastUpdateTimestamp int                    `json:"lastUpdateTimestamp"`
	Layer               int                    `json:"layer"`
	Name                string                 `json:"name"`
	Properties          map[string]interface{} `json:"properties"`
	State               map[string]interface{} `json:"state"`
	SyncedElems         []SyncElem             `json:"synced"`
	SyncedData          map[string][]SyncData  `json:"synchronizationData"`

	Tags []string `json:"tags"`
}

// SyncData returned in a TopologyStream Query
type SyncData struct {
	Data             map[string]interface{} `json:"data"`
	SourceProperties map[string]interface{} `json:"sourceProperties"`
}

// SyncElem return in a Topology Query with full component
type SyncElem struct {
	Type               string             `json:"_type"`
	ExtTopologyElement ExtTopologyElement `json:"extTopologyElement"`
}

type ExtTopologyElement struct {
	ElementTypeTag   string                 `json:"elementTypeTag"`
	ExternalId       string                 `json:"externalId"`
	Id               int64                  `json:"id"`
	Identifiers      []string               `json:"identifiers"`
	Data             map[string]interface{} `json:"data"`
	Tags             map[string]interface{} `json:"tags"`
	SourceProperties map[string]interface{} `json:"sourceProperties"`
}
type Label string

func (t *Label) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	var label map[string]string
	if err := json.Unmarshal(data, &label); err != nil {
		return err
	}

	*t = Label(label["name"])
	return nil
}

type TopoQueryResponse struct {
	Success bool            `json:"success"`
	Errors  []*ErrorMsg     `json:"errors"`
	Data    []SyncComponent `json:"data"`
}

type ErrorResp struct {
	Errors []*ErrorMsg `json:"errors"`
}

type SuccessResp struct {
	Result []SyncComponent `json:"result"`
}

type ErrorMsg struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
}

type scriptRequest struct {
	ReqType string `json:"_type"`
	Body    string `json:"body"`
}

type ViewSnapshotResponse struct {
	Components []ViewComponent
}

type ViewComponent struct {
	ID                  int64   `json:"id"`
	Name                string  `json:"name"`
	Description         string  `json:"description"`
	LastUpdateTimestamp int64   `json:"lastUpdateTimestamp"`
	Type                int64   `json:"type"`
	Layer               int     `json:"layer"`
	Domain              int     `json:"domain"`
	Environments        []int64 `json:"environments"`
	State               struct {
		ID                    int64  `json:"id"`
		LastUpdateTimestamp   int64  `json:"lastUpdateTimestamp"`
		HealthState           string `json:"healthState"`
		PropagatedHealthState string `json:"propagatedHealthState"`
		Type                  string `json:"_type"`
	} `json:"state"`
	OutgoingRelations []int64           `json:"outgoingRelations"`
	IncomingRelations []int64           `json:"incomingRelations"`
	Synchronized      bool              `json:"synchronized"`
	FailingChecks     []any             `json:"failingChecks"`
	RetrievalSource   string            `json:"retrievalSource"`
	Identifiers       []string          `json:"identifiers"`
	Tags              []string          `json:"tags"`
	Properties        map[string]string `json:"properties"`
	InternalType      string            `json:"_type"`
}

type ViewSnapshotRequest struct {
	Type         string               `json:"_type"`
	Metadata     ViewSnapshotMetadata `json:"metadata"`
	Query        string               `json:"query"`
	QueryVersion string               `json:"queryVersion"`
}

type ViewSnapshotMetadata struct {
	Type                  string `json:"_type"`
	GroupingEnabled       bool   `json:"groupingEnabled"`
	ShowIndirectRelations bool   `json:"showIndirectRelations"`
	MinGroupSize          int    `json:"minGroupSize"`
	GroupedByLayer        bool   `json:"groupedByLayer"`
	GroupedByDomain       bool   `json:"groupedByDomain"`
	GroupedByRelation     bool   `json:"groupedByRelation"`
	ShowCause             string `json:"showCause"`
	AutoGrouping          bool   `json:"autoGrouping"`
	ConnectedComponents   bool   `json:"connectedComponents"`
	NeighboringComponents bool   `json:"neighboringComponents"`
	QueryTime             int64  `json:"queryTime"`
}

func newViewSnapshotRequest() *ViewSnapshotRequest {
	return &ViewSnapshotRequest{
		Type: "ViewSnapshotRequest",
		Metadata: ViewSnapshotMetadata{
			Type:         "QueryMetadata",
			MinGroupSize: 2,
			ShowCause:    "NONE",
			QueryTime:    time.Now().Unix(),
		},
		QueryVersion: "0.0.1",
	}
}

type NodeType struct {
	TypeName            string `json:"typeName"`
	ID                  int64  `json:"id"`
	LastUpdateTimestamp int64  `json:"lastUpdateTimestamp"`
	Identifier          string `json:"identifier"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	OwnedBy             string `json:"ownedBy"`
	Manual              bool   `json:"manual"`
	IsSettingsNode      bool   `json:"isSettingsNode"`
	Type                string `json:"_type"`
}
