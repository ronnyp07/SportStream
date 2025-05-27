package cccopy

import (
	"encoding/json"
	"fmt"
)

func DeepCopy(src interface{}, dest interface{}) (err error) {
	srcBinary, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("unable to marshal %v to %v, source: %v", src, dest, err)
	}

	if err := json.Unmarshal(srcBinary, &dest); err != nil {
		return fmt.Errorf("unable to unmarshal %v to %v, source: %v", src, dest, err)
	}

	return err
}

func CopySlice[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

func CopyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
