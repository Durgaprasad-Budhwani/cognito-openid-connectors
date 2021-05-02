package common

import (
	"github.com/pkg/errors"
)

func GetStringKey(m map[string]string, key string) (string, error) {
	item, ok := m[key]
	if ok {
		return item, nil
	}
	return "", errors.New(key + "is required")
}
