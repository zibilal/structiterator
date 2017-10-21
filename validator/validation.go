package validator

import (
	"fmt"
	"errors"
)

type Validation struct {}

func (v Validation) Required(value interface{}, key string, defaultError string) error {
	if IsEmpty(value) {
		if defaultError == "" {
			return fmt.Errorf("%s is required", key)
		}
		return errors.New(defaultError)
	} else {
		return nil
	}
}

func (v Validation) After(structValue interface{}, key1, key2 string, defaultError string) error {
	return nil
}

func (v Validation) Date(value interface{}, key string, defaultError string) error {
	return nil
}

func (v Validation) Phone(value interface{}, key string, defaultError string) error {
	return nil
}

func (v Validation) Email(value interface{}, key string, defaultError string) error {
	return nil
}
