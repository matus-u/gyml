package gyml

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrRootNodeNotSet     = errors.New("rootNode not set")
	ErrKeyNotFound        = errors.New("key not found")
	ErrEmptyDocumentNode  = errors.New("empty document node provided")
	ErrUnexpectedNodeKind = errors.New("unexpected node kind provided")
	ErrInvalidIndexFormat = errors.New("invalid index format")
	ErrIndexOutOfBound    = errors.New("provided index out of bound")
	ErrInvalidKeysList    = errors.New("invalid keys list")
)

func DeleteValue(root *yaml.Node, keys ...string) error {

	if len(keys) == 0 {
		return ErrInvalidKeysList
	}

	if root == nil {
		return ErrRootNodeNotSet
	}
	return deleteValue(root, keys...)
}

// Returns values on the path defined by list of keys
// Examples:
// GetValue[int](&number, "persons_list", "[10]", "age") - get age property of 10th person in person_list, deserialize to *int
func GetValue[DataType any](rootNode *yaml.Node, keys ...string) (*DataType, error) {

	if rootNode == nil {
		return nil, ErrRootNodeNotSet
	}

	node, err := getValue(rootNode, keys...)
	if err != nil {
		return nil, err
	}

	var value DataType
	if err := node.Decode(&value); err != nil {
		return nil, fmt.Errorf("GetValue: cannot unmarshall yaml node value: %w", err)
	}
	normalizeEmptySlice(&value)
	return &value, nil

}

func parseValidIndex(indexStr string) (int, error) {
	if len(indexStr) < 3 {
		return 0, ErrInvalidIndexFormat
	}

	if indexStr[0] != '[' || indexStr[len(indexStr)-1] != ']' {
		return 0, ErrInvalidIndexFormat
	}

	index, err := strconv.Atoi(indexStr[1 : len(indexStr)-1])

	if err != nil {
		return 0, ErrInvalidIndexFormat
	}

	return index, nil
}

func getValue(node *yaml.Node, keys ...string) (*yaml.Node, error) {

	// final recursion
	if len(keys) == 0 {
		return node, nil
	}

	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return nil, ErrEmptyDocumentNode
		}
		return getValue(node.Content[0], keys...)
	}

	if node.Kind == yaml.SequenceNode {
		index, err := parseValidIndex(keys[0])
		if err != nil {
			return nil, err
		}

		if index < 0 || index >= len(node.Content) {
			return nil, ErrIndexOutOfBound
		}

		return getValue(node.Content[index], keys[1:]...)
	}

	if node.Kind == yaml.MappingNode {
		// Content is sorted as key1,value1,key2,value2...
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == keys[0] {
				return getValue(node.Content[i+1], keys[1:]...)
			}
		}
		return nil, fmt.Errorf("%w: %s", ErrKeyNotFound, keys[0])
	}

	return nil, fmt.Errorf("%w: key: %s", ErrUnexpectedNodeKind, keys[0])
}

func normalizeEmptySlice[T any](v *T) {
	if v == nil {
		return
	}

	rv := reflect.ValueOf(v).Elem()
	if rv.Kind() == reflect.Slice && rv.IsNil() {
		rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
	}
}

func deleteValue(node *yaml.Node, keys ...string) error {

	if len(keys) == 0 || node == nil {
		return ErrInvalidKeysList
	}

	if node.Kind == yaml.DocumentNode {
		if len(node.Content) > 0 {
			return deleteValue(node.Content[0], keys...)
		}
		return ErrEmptyDocumentNode
	}

	if node.Kind == yaml.SequenceNode {
		return nil
	}

	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == keys[0] {
				if len(keys) == 1 {
					node.Content = slices.Delete(node.Content, i, i+2)
					return nil
				}
				valueNode := node.Content[i+1]
				retVal := deleteValue(valueNode, keys[1:]...)
				if retVal == nil && isEmptyNode(valueNode) {
					node.Content = slices.Delete(node.Content, i, i+2)
				}
				return retVal
			}
		}
		return fmt.Errorf("%w: %s", ErrKeyNotFound, keys[0])
	}

	if node.Kind == yaml.ScalarNode {
		return fmt.Errorf("%w: remaining path: %s", ErrInvalidKeysList, strings.Join(keys, "."))
	}

	return fmt.Errorf("%w: key: %s", ErrUnexpectedNodeKind, keys[0])
}

func isEmptyNode(node *yaml.Node) bool {
	if node.Kind == yaml.MappingNode || node.Kind == yaml.SequenceNode {
		return len(node.Content) == 0
	}
	return false
}
