package gyml

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

const testYAML = `
clients:
  - name: first_client
    surname: first_surname
  - name: second_client
    surname: second_surname
servers:
  server1:
    host: server1.local
    port: 9001
  server2:
    host: server2.local
    port: 9002
ints:
  - 10
  - 20
  - 30
`

const emptyYAML = ``
const listYAML = `
- 10
- 20
`

func TestEmptyRootNode(t *testing.T) {
	value, err := GetValue[int](nil, "Unknown", "[10]")
	require.Nil(t, value)
	require.Equal(t, ErrRootNodeNotSet, err)
}

func TestGetValue(t *testing.T) {
	var root yaml.Node
	var rootList yaml.Node
	var rootEmpty yaml.Node

	err := yaml.Unmarshal([]byte(testYAML), &root)
	require.NoError(t, err)

	err = yaml.Unmarshal([]byte(emptyYAML), &rootEmpty)
	require.NoError(t, err)

	err = yaml.Unmarshal([]byte(listYAML), &rootList)
	require.NoError(t, err)

	host1, err := GetValue[string](&root, "servers", "server1", "host")
	require.NoError(t, err)
	require.Equal(t, "server1.local", *host1)

	second_surname, err := GetValue[string](&root, "clients", "[1]", "surname")
	require.NoError(t, err)
	require.Equal(t, "second_surname", *second_surname)

	ints, err := GetValue[[]int](&root, "ints")
	require.NoError(t, err)
	require.Equal(t, []int{10, 20, 30}, *ints)

	ints, err = GetValue[[]int](&root, "non_existent_ints")
	require.ErrorIs(t, err, ErrKeyNotFound)
	require.Nil(t, ints)

	ints, err = GetValue[[]int](&rootList)
	require.NoError(t, err)
	require.Equal(t, []int{10, 20}, *ints)

	ints, err = GetValue[[]int](&rootList, "[*]")
	require.Equal(t, ErrInvalidIndexFormat, err)
	require.Nil(t, ints)

	val, err := GetValue[int](&rootList, "[1]")
	require.NoError(t, err)
	require.Equal(t, *val, 20)

	val, err = GetValue[int](&rootList, "[0]")
	require.NoError(t, err)
	require.Equal(t, *val, 10)

	val, err = GetValue[int](&rootList, "[-1]")
	require.Equal(t, ErrIndexOutOfBound, err)
	require.Nil(t, val)

	val, err = GetValue[int](&rootList, "[3]")
	require.Equal(t, ErrIndexOutOfBound, err)
	require.Nil(t, val)

	val, err = GetValue[int](&rootList, "[25]")
	require.Equal(t, ErrIndexOutOfBound, err)
	require.Nil(t, val)

	ints, err = GetValue[[]int](&rootEmpty)
	require.NoError(t, err)
	require.Equal(t, []int{}, *ints)
}

func TestDeleteValue(t *testing.T) {
	var root yaml.Node
	var rootList yaml.Node
	var rootEmpty yaml.Node

	err := yaml.Unmarshal([]byte(testYAML), &root)
	require.NoError(t, err)

	err = yaml.Unmarshal([]byte(emptyYAML), &rootEmpty)
	require.NoError(t, err)

	err = yaml.Unmarshal([]byte(listYAML), &rootList)
	require.NoError(t, err)

	err = DeleteValue(&rootEmpty, "servers", "server1", "host")
	require.ErrorIs(t, err, ErrUnexpectedNodeKind)

	err = DeleteValue(&root, "servers", "server1", "host")
	require.NoError(t, err)

	err = DeleteValue(&root, "servers", "server1", "host")
	require.ErrorIs(t, err, ErrKeyNotFound)

	_, err = GetValue[string](&root, "servers", "server1", "host")
	require.ErrorIs(t, err, ErrKeyNotFound)

	err = DeleteValue(&root, "servers", "server1", "port")
	require.NoError(t, err)

	err = DeleteValue(&root, "ints", "[1]")
	require.NoError(t, err)

	val, err := GetValue[int](&root, "ints", "[1]")
	require.NoError(t, err)
	require.Equal(t, *val, 30)

	err = DeleteValue(&root, "ints", "[-25]")
	require.Equal(t, ErrIndexOutOfBound, err)

	err = DeleteValue(&root, "ints", "[25]")
	require.Equal(t, ErrIndexOutOfBound, err)

	err = DeleteValue(&root, "ints")
	require.NoError(t, err)

	err = DeleteValue(&rootList, "[1]")
	require.NoError(t, err)

	val, err = GetValue[int](&rootList, "[0]")
	require.NoError(t, err)
	require.Equal(t, *val, 10)

	err = DeleteValue(&rootList, "[1]")
	require.Equal(t, ErrIndexOutOfBound, err)

}
