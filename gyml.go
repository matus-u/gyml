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
	ErrRootNodeNotSet        = errors.New("rootNode not set")
	ErrKeyNotFound           = errors.New("key not found")
	ErrEmptyDocumentNode     = errors.New("empty document node provided")
	ErrUnexpectedNodeKind    = errors.New("unexpected node kind provided")
	ErrInvalidIndexFormat    = errors.New("invalid index format")
	ErrIndexOutOfBound       = errors.New("provided index out of bound")
	ErrInvalidKeysList       = errors.New("invalid keys list")
	ErrScalarSetAttempt      = errors.New("cannot iterate over scalar node")
	ErrCannotEncodeData      = errors.New("cannot marshall provided data to yaml")
	ErrAppendInWrongPosition = errors.New("cannot append in the middle of the keys")
)

// SetValue is used to set values on a key's path within a YAML node.
// If a part of the path does not exist, it is created with the correct type
// based on the key (e.g., "[]" for a sequence, "name" for a map).
//
// root: The root YAML node to modify.
// data: The data to be inserted, which can be any type serializable into YAML
//
//	(ScalarType, MapType, SeqType).
//
// keys: A list of keys defining the path in the YAML structure where the value
//
//	should be set.
//
// Examples:
//
//	SetValue(&root, TestSomeStruct{Name: "Adam", Age: 30}, "Company", "CEO")
//	  - Sets TestSomeStruct at the path /Company/CEO.
//	SetValue(&root, 35, "some_list", "[]")
//	  - Appends the item 35 to the "some_list" sequence.
//	SetValue(&root, 12, "some_list", "[8]")
//	  - Sets the value 12 at index 8 in "some_list" (with range checking).
//	SetValue(&root, "Matus", "Company", "CEO", "Name")
//	  - Sets the scalar value "Matus" at /Company/CEO/Name.
//	SetValue(&root, true, "settings", "notifications")
//	  - Sets a boolean value 'true' at /settings/notifications.
//	SetValue(&root, 3.14, "math", "pi")
//	  - Sets a float value '3.14' at /math/pi.
//	SetValue(&root, map[string]string{"city": "London", "country": "UK"}, "user", "address")
//	  - Sets a nested map at /user/address.
func SetValue[DataType any](root *yaml.Node, data DataType, keys ...string) error {
	if root == nil {
		return ErrRootNodeNotSet
	}

	return setValue(root, data, keys...)
}

// DeleteValue removes a value from a YAML node at the specified path.
//
// root: The root YAML node from which to delete the value.
// keys: A list of keys defining the path to the value to be deleted.
func DeleteValue(root *yaml.Node, keys ...string) error {

	if len(keys) == 0 {
		return ErrInvalidKeysList
	}

	if root == nil {
		return ErrRootNodeNotSet
	}
	return deleteValue(root, keys...)
}

// GetValue retrieves a value from a YAML node at the path defined by a list of keys.
// The retrieved value is deserialized into the specified DataType.
//
// rootNode: The root YAML node from which to retrieve the value.
// keys:     A list of keys defining the path to the desired value.
//
// Examples:
//
//	GetValue[int](&number, "persons_list", "[10]", "age")
//	  - Retrieves the "age" property of the 10th person in "persons_list"
//	    and deserializes it into an *int.
//	GetValue[bool](&notificationsEnabled, "settings", "notifications")
//	  - Retrieves a boolean value from /settings/notifications.
//	GetValue[float64](&piValue, "math", "pi")
//	  - Retrieves a float64 value from /math/pi.
//	GetValue[string](&city, "user", "address", "city")
//	  - Retrieves a string value from /user/address/city.
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
		return nil, fmt.Errorf("GetValue: cannot decode yaml node value: %w", err)
	}
	normalizeEmptySlice(&value)
	return &value, nil

}

func parseValidIndex(indexStr string, node *yaml.Node) (int, error) {
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

	if index < 0 || index >= len(node.Content) {
		return 0, ErrIndexOutOfBound
	}

	return index, nil
}

func setValue[DataType any](node *yaml.Node, data DataType, keys ...string) error {

	if len(keys) == 0 {
		contentNode, err := createContentNode(data)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCannotEncodeData, err)
		}
		*node = *contentNode
		return nil
	}

	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == keys[0] {
				return setValue(node.Content[i+1], data, keys[1:]...)
			}
		}

		return appendDataToContent(node, data, keys...)
	}

	if node.Kind == yaml.SequenceNode {

		if keys[0] == "[]" {
			if len(keys) == 1 {
				return appendDataToContent(node, data, keys...)
			}
			return fmt.Errorf("%w: following key: %s", ErrAppendInWrongPosition, keys[1])
		}

		index, err := parseValidIndex(keys[0], node)
		if err != nil {
			return err
		}

		return setValue(node.Content[index], data, keys[1:]...)
	}

	if node.Kind == yaml.ScalarNode {
		return fmt.Errorf("%w: %s", ErrScalarSetAttempt, keys[0])
	}

	if node.Kind == yaml.DocumentNode {
		if len(node.Content) > 0 {
			return setValue(node.Content[0], data, keys...)
		}

		return appendDataToContent(node, data, keys...)
	}

	return fmt.Errorf("%w: key: %s", ErrUnexpectedNodeKind, keys[0])
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
		index, err := parseValidIndex(keys[0], node)
		if err != nil {
			return nil, err
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

		index, err := parseValidIndex(keys[0], node)
		if err != nil {
			return err
		}

		if len(keys) == 1 {
			node.Content = slices.Delete(node.Content, index, index+1)
			return nil
		}

		retVal := deleteValue(node.Content[index], keys[1:]...)
		// delete empty list itself when it is empty after deleting my last child
		if retVal == nil && isEmptyNode(node.Content[index]) {
			node.Content = slices.Delete(node.Content, index, index+1)
		}

		return retVal
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

				// delete empty map itself when it is empty after deleting my last child
				if retVal == nil && isEmptyNode(valueNode) {
					node.Content = slices.Delete(node.Content, i, i+2)
				}
				return retVal
			}
		}
		return fmt.Errorf("%w: %s", ErrKeyNotFound, keys[0])
	}

	if node.Kind == yaml.ScalarNode {
		return fmt.Errorf("%w: unresolved path: %s", ErrInvalidKeysList, strings.Join(keys, "."))
	}

	return fmt.Errorf("%w: key: %s", ErrUnexpectedNodeKind, keys[0])
}

func isEmptyNode(node *yaml.Node) bool {
	if node.Kind == yaml.MappingNode || node.Kind == yaml.SequenceNode {
		return len(node.Content) == 0
	}
	return false
}

// createTypedEnvelope recursively wraps the provided data into a nested structure
// of maps or slices based on the provided restKeys.
// If a key is "[]", it wraps the result in a slice.
// Otherwise, it wraps the result in a map with the key as the property name.
func createTypedEnvelope[DataType any](data DataType, restKeys ...string) any {
	if len(restKeys) == 0 {
		// Base case: no more keys to process, return the raw data
		return data
	}

	if restKeys[0] == "[]" {
		// List envelope: wrap the next level in a slice
		return []any{createTypedEnvelope(data, restKeys[1:]...)}
	}

	// Map envelope: wrap the next level in a map using the current key
	return map[string]any{restKeys[0]: createTypedEnvelope(data, restKeys[1:]...)}
}

// wrap any data in ContentNode to set/append to another Node
func createContentNode[DataType any](data DataType) (*yaml.Node, error) {
	node := yaml.Node{}
	if err := node.Encode(data); err != nil {
		return nil, fmt.Errorf("cannot encode value to yaml node: %w", err)
	}

	return &node, nil
}

func appendDataToContent[DataType any](node *yaml.Node, data DataType, keys ...string) error {
	contentNode, err := createContentNode(createTypedEnvelope(data, keys...))
	if err != nil {
		return err
	}
	node.Content = append(node.Content, contentNode.Content...)
	return nil
}
