package gyml

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

const testYAML = `
servers:
  server1:
    host: server1.local
    port: 9001
  server2:
    host: server2.local
    port: 9002
`

func TestEmptyRootNode(t *testing.T) {
	value, err := GetValue[int](nil, "Unknown", "[10]")
	require.Nil(t, value)
	require.Equal(t, ErrRootNodeNotSet, err)
}

func TestGetValue(t *testing.T) {
	var root yaml.Node

	err := yaml.Unmarshal([]byte(testYAML), &root)
	require.NoError(t, err)

	host1, err := GetValue[string](&root, "servers", "server1", "host")
	require.NoError(t, err)
	require.Equal(t, *host1, "server1.local")

}
