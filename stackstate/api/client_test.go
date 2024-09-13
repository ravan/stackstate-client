package api

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	sts "github.com/ravan/stackstate-client/stackstate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

var opts = &slog.HandlerOptions{Level: slog.LevelDebug}
var handler = slog.NewJSONHandler(os.Stdout, opts)
var logger = slog.New(handler)

func init() { slog.SetDefault(logger) }

func TestClientConnection(t *testing.T) {
	conf := getConfig(t)
	client := NewClient(conf)
	status, err := client.Status()
	require.NoError(t, err, `Not expecting err %v`, err)
	assert.Equal(t, status.Version.Major, 6)
}

func TestTopologyQuery(t *testing.T) {
	conf := getConfig(t)
	client := NewClient(conf)
	res, err := client.TopologyQuery("type = 'service' and label in ('namespace:kube-system')", "", false)
	require.NoError(t, err)
	require.True(t, res.Success, `Expected to be successful but was %s`, toJson(res))
	assert.True(t, len(res.Data) > 0)
	fmt.Println(toJson(res))
}

func TestTopologyStreamQuery(t *testing.T) {
	conf := getConfig(t)
	client := NewClient(conf)
	res, err := client.TopologyStreamQuery("type = 'service' and label in ('namespace:kube-system')", "", true)
	require.NoError(t, err)
	require.True(t, res.Success, `Expected to be successful but was %s`, toJson(res))
	assert.True(t, len(res.Data) > 0)
	fmt.Println(toJson(res))
}

func toJson(a any) string {
	marshal, err := json.Marshal(a)
	if err != nil {
		fmt.Printf("Failed to marshall json. %v", err)
	}
	return string(marshal)
}

func getConfig(t *testing.T) *sts.StackState {
	require.NoError(t, godotenv.Load())
	return &sts.StackState{
		ApiUrl:   os.Getenv("STS_URL"),
		ApiKey:   os.Getenv("STS_API_KEY"),
		ApiToken: os.Getenv("STS_TOKEN"),
	}
}
