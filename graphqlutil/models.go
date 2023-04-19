package graphqlutil

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
)

type ElementList[T comparable] []T

func (list ElementList[T]) Merge(other ElementList[T]) ElementList[T] {
	uniqueElement := make(map[T]struct{})
	merged := make(ElementList[T], 0)
	for _, item := range list {
		if _, ok := uniqueElement[item]; ok {
			continue
		}
		merged = append(merged, item)
		uniqueElement[item] = struct{}{}
	}
	for _, item := range other {
		if _, ok := uniqueElement[item]; ok {
			continue
		}
		merged = append(merged, item)
		uniqueElement[item] = struct{}{}
	}
	return merged
}

func (list ElementList[T]) ToBytes() []byte {
	if list == nil {
		return []byte(`[]`)
	}
	ret, _ := json.Marshal(list)
	return ret
}

func (list ElementList[T]) HasDuplicate() bool {
	uniqueElement := make(map[T]struct{})
	for _, item := range list {
		if _, ok := uniqueElement[item]; ok {
			return true
		}
		uniqueElement[item] = struct{}{}
	}
	return false
}

func (list ElementList[T]) Value() (driver.Value, error) {
	return json.Marshal(list)
}

func (list *ElementList[T]) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), list)
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (list ElementList[T]) UnmarshalGQL(v interface{}) error {

	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("list [%v] must be a string", v)
	}
	return json.Unmarshal([]byte(s), &list)
}

// MarshalGQL implements the graphql.Marshaler interface
func (list ElementList[T]) MarshalGQL(w io.Writer) {
	encoder := json.NewEncoder(w)
	encoder.Encode(list)
}
