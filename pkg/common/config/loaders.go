package config

import (
	"fmt"
	"github.com/Nevermore12321/ShIM/pkg/jsonx"
	"github.com/Nevermore12321/ShIM/pkg/mapping"
	"reflect"
	"strings"
)

const (
	JsonTagName      = "json" // json tag
	JsonTagSeparator = ','    // json tag separator
)

/*
	field information for configuration file
	e.g.1
		type a struct {
			b struct{}
			c struct{}
		}
	e.g.2
		type a struct {
			b []struct{}
		}

	notes:
	1. children and mapField should not be both filled.
	2. named fields and map cannot be bound to the same field name.
*/

type fieldInfo struct {
	mapField *fieldInfo            // last layer
	children map[string]*fieldInfo // has children
}

type conflictKeyError struct {
	key string
}

func (e conflictKeyError) Error() string {
	return fmt.Sprintf("conflict key %s, pay attention to anonymous fields", e.key)
}

func newConflictKeyError(key string) conflictKeyError {
	return conflictKeyError{key: key}
}

// LoadFromJsonBytes load config into v object from json bytes
func LoadFromJsonBytes(content []byte, v any) error {
	// all field info, is a filedInfo object
	info, err := buildFieldsInfo(reflect.TypeOf(v), "")
	if err != nil {
		return nil
	}
	fmt.Println(info)

	// transfer content to a map
	var m map[string]any
	err = jsonx.Unmarshal(content, &m)
	if err != nil {
		return nil
	}

	lowerCaseKeyMap := toLowerCaseKeyMap(m, info)

	fmt.Println(lowerCaseKeyMap)
	return err
}

// LoadFromYamlBytes load config into v object from yaml bytes
func LoadFromYamlBytes(content []byte, v any) error {
	return nil
}

// Perform different operations based on the type being converted
// param rt: Type being converted
// param fieldName: configuration filed name full path: e.g. a.b.c
// return fullName: converted type(any type) to fieldInfo type
func buildFieldsInfo(fieldType reflect.Type, fullName string) (*fieldInfo, error) {
	// if rt is a point, Obtain the type of the converted object
	fieldType = mapping.Dereference(fieldType)

	// Determine the type of the converted object
	switch fieldType.Kind() {
	case reflect.Struct:
		// converted type is struct, continue to build
		// e.g. type a struct { b struct, c struct}
		return buildStructFiledInfo(fieldType, fullName)
	case reflect.Slice, reflect.Array:
		// converted type is Slice, continue to build
		// Elem() returns a type's element type and will panic if the type's Kind is not Array, Chan, Map, Pointer, or Slice.
		// e.g. type a struct { b []struct }
		return buildFieldsInfo(mapping.Dereference(fieldType.Elem()), fullName)
	case reflect.Chan, reflect.Func:
		return nil, fmt.Errorf("unsupported type: %s", fieldType.Kind())
	default:
		// converted type maybe is map or other type
		return &fieldInfo{
			children: make(map[string]*fieldInfo),
		}, nil
	}
}

// converted type is struct
// fieldType is struct type
func buildStructFiledInfo(fieldType reflect.Type, fullName string) (*fieldInfo, error) {
	info := &fieldInfo{
		children: make(map[string]*fieldInfo),
	}

	// NumField returns a struct type's field count.
	// Traverse every field
	for i := 0; i < fieldType.NumField(); i++ {
		// obtain a struct type's i'th field.
		field := fieldType.Field(i)
		// Can other packages import this feed. Is the field uppercase or lowercase
		// If not visible outside the package, skip
		if !field.IsExported() {
			continue
		}

		// if converted struct object's value is a pointer, Get the data pointed to by the pointer
		// e.g. type A struct { b: pointer }
		pt := mapping.Dereference(field.Type)

		// obtain json tag name. That is, the field name in the configuration file
		childName := getTagName(field)
		childLowerName := strings.ToLower(childName)

		// parent.child
		childFullName := getFullName(fullName, childLowerName)

		// if current field is an embedded field(anonymous fields)
		// e.g. type B struct{},  type A struct { B }
		if field.Anonymous {
			err := buildAnonymousFieldInfo(info, childLowerName, pt, childFullName)
			if err != nil {
				return nil, err
			}

		} else { // named field. e.g. type B struct{},  type A struct { b: B }
			err := buildNamedFieldInfo(info, childLowerName, pt, childFullName)
			if err != nil {
				return nil, err
			}
		}
	}
	return info, nil
}

// getTagName get the tag name of the given field, if no tag name, use file.Name.
// returned on tags like `json:""` and `json:",optional"`.
func getTagName(field reflect.StructField) string {
	// find tag of the given tag name, e.g. `json:""` tag name is json
	tagValue, ok := field.Tag.Lookup(JsonTagName)
	if !ok {
		// if no tag, return field name
		return field.Name
	}

	// `json:"test,optional"` take out the test part. The field name defined by the tag
	index := strings.IndexByte(tagValue, JsonTagSeparator)
	if index >= 0 {
		tagValue = tagValue[:index]
	}
	tagValue = strings.TrimSpace(tagValue)
	if len(tagValue) > 0 {
		return tagValue
	}
	return field.Name
}

// e.g. type A struct { b: B }
func buildNamedFieldInfo(info *fieldInfo, lowerCaseName string, fieldType reflect.Type, fullName string) error {
	var (
		childFieldInfo *fieldInfo
		err            error
	)

	switch fieldType.Kind() {
	case reflect.Struct:
		childFieldInfo, err = buildFieldsInfo(fieldType, fullName)
		if err != nil {
			return err
		}
	case reflect.Slice, reflect.Array:
		childFieldInfo, err = buildFieldsInfo(fieldType.Elem(), fullName)
		if err != nil {
			return err
		}
	case reflect.Map:
		// converted type is map, it is a pointer, need obtain pointed data.
		// Recursively obtain all elements of map type
		// last layer
		elementInfo, err := buildFieldsInfo(mapping.Dereference(fieldType.Elem()), fullName)
		if err != nil {
			return err
		}
		childFieldInfo = &fieldInfo{
			children: make(map[string]*fieldInfo),
			mapField: elementInfo,
		}
	default:
		childFieldInfo, err = buildFieldsInfo(fieldType, fullName)
		if err != nil {
			return err
		}
	}
	return addOrMergeFieldInfo(info, lowerCaseName, childFieldInfo, fullName)
}

// e.g. type A struct { B }
func buildAnonymousFieldInfo(info *fieldInfo, lowerCaseName string, fieldType reflect.Type, fullName string) error {
	switch fieldType.Kind() {
	case reflect.Struct:
		// e.g. type A struct { B }, continue to build A.B
		// To the anonymous attribute, it is already the last layer
		fields, err := buildFieldsInfo(fieldType, fullName)
		if err != nil {
			return err
		}
		// add or merge info
		for k, v := range fields.children {
			err := addOrMergeFieldInfo(info, k, v, fullName)
			if err != nil {
				return err
			}
		}

	case reflect.Map:
		// converted type is map, it is a pointer, need obtain pointed data.
		// Recursively obtain all elements of map type
		elementInfo, err := buildFieldsInfo(mapping.Dereference(fieldType.Elem()), fullName)
		if err != nil {
			return err
		}

		// If the same key already exists
		if _, ok := info.children[lowerCaseName]; ok {
			return newConflictKeyError(lowerCaseName)
		}

		// Otherwise, add the current field to the info
		info.children[lowerCaseName] = &fieldInfo{
			children: make(map[string]*fieldInfo),
			mapField: elementInfo,
		}

	default:
		// If the same key already exists
		if _, ok := info.children[lowerCaseName]; ok {
			return newConflictKeyError(lowerCaseName)
		}

		// Otherwise, add the current field to the info
		info.children[lowerCaseName] = &fieldInfo{
			children: make(map[string]*fieldInfo),
		}
	}
	return nil
}

// obtain full name fieldNmae.childName
// e.g. type A struct { b: B }, fullNmae is a.b
func getFullName(fieldName, childName string) string {
	if fieldName == "" {
		return childName
	}

	return strings.Join([]string{fieldName, childName}, ".")
}

func addOrMergeFieldInfo(info *fieldInfo, key string, child *fieldInfo, fullName string) error {
	// if fieldInfo has both children and mapField, key conflict
	if prev, ok := info.children[key]; ok {
		if child.mapField != nil {
			return newConflictKeyError(key)
		}
		// already has this key，merge
		if err := mergeFields(prev, key, child.children, fullName); err != nil {
			return err
		}
	} else { // has not this key, add
		info.children[key] = child
	}
	return nil
}

func mergeFields(prev *fieldInfo, key string, children map[string]*fieldInfo, fullName string) error {
	if len(prev.children) == 0 || len(children) == 0 {
		return newConflictKeyError(key)
	}

	// merge all fileds
	for k, v := range children {
		if _, ok := prev.children[k]; ok {
			return newConflictKeyError(k)
		}

		prev.children[k] = v
	}

	return nil
}

// Convert the unmarshal key to lowercase characters
func toLowerCaseKeyMap(m map[string]any, info *fieldInfo) any {
	result := make(map[string]any)

	for k, v := range m {
		infoChildValue, ok := info.children[k]
		if ok {
			result[k] = toLowerCaseInterface(v, infoChildValue)
			continue
		}

		lowerKey := strings.ToLower(k)
		infoChildValue, ok = info.children[lowerKey]
		if ok {
			result[lowerKey] = toLowerCaseInterface(v, infoChildValue)
		} else if info.mapField != nil {
			result[k] = toLowerCaseInterface(v, info.mapField)
		} else {
			result[k] = v
		}
	}
	return result
}

func toLowerCaseInterface(v any, info *fieldInfo) any {
	switch vv := v.(type) {
	case map[string]any: // if it is not the last layer
		return toLowerCaseKeyMap(vv, info)
	case []any:
		var array []any
		for _, vvv := range vv {
			array = append(array, toLowerCaseInterface(vvv, info))
		}
		return array
	default:
		return v
	}
}
