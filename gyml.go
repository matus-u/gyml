package gyml

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	ErrRootNodeNotSet     = errors.New("error: rootNode not set")
	ErrKeyNotFound        = errors.New("error: key not found")
	ErrEmptyDocumentNode  = errors.New("error: empty document node provided")
	ErrUnexpectedNodeKind = errors.New("error: unexpected node kind provided")
)

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
	return &value, nil

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

	if node.Kind == yaml.MappingNode {
		// Content is sorted as key1,value1,key2,value2...
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == keys[0] {
				return getValue(node.Content[i+1], keys[1:]...)
			}
		}
	}

	return nil, ErrUnexpectedNodeKind
}
