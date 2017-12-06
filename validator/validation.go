package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"
	"strings"
)

type Validation struct {
	PhoneFormat string
	EmailFormat string
	DateFormat  string
}

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

func (v Validation) CondRequired(structValue interface{}, key string,  zeValue interface{}, keyCompare, valueCompare string, defaultError string) error {

	val := reflect.ValueOf(structValue)
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("Bad value, expected struct value, got ")
	}

	for i:=0; i < val.NumField(); i++ {
		fv := val.Field(i)
		ft := typ.Field(i)

		vkey := ft.Tag.Get("vkey")
		if (vkey != "" && vkey == keyCompare) || (vkey == "" && ft.Name == keyCompare) {
			zVal := fmt.Sprintf("%v", reflect.Indirect(fv))
			split := strings.Split(valueCompare, "|")
			for _, sp := range split {
				if zVal == sp {
					if IsEmpty(zeValue) {
						if defaultError == "" {
							return fmt.Errorf("%s is required", key)
						} else {
							return errors.New(defaultError)
						}
					}
				}
			}
		}

	}
	return nil
}

func (v Validation) AfterDate(structValue interface{}, key1, key2 string, defaultError string) error {

	val1, val2, err := v.after(structValue, key1, key2, defaultError)

	if err != nil {
		return err
	}

	ival1, found := val1.Interface().(string)
	if !found {
		return fmt.Errorf("Expected type string got %s", val1.Type())
	}

	ival2, found := val2.Interface().(string)
	if !found {
		return fmt.Errorf("Expected type string got %s", val2.Type())
	}

	if ival1 == "" || ival2 == "" {
		return fmt.Errorf("Expected both values not empty")
	}

	var dateLayout string
	if v.DateFormat == "" {
		dateLayout = "01/02/2006"
	} else {
		dateLayout = v.DateFormat
	}

	time1, err := time.Parse(dateLayout, ival1)
	if err != nil {
		return err
	}

	time2, err := time.Parse(dateLayout, ival2)
	if err != nil {
		return err
	}

	if !time1.After(time2) {
		if defaultError == "" {
			return fmt.Errorf("Invalid %s should be after %s", key1, key2)
		}

		return errors.New(defaultError)
	}

	return nil
}

func (v Validation) after(structValue interface{}, key1, key2 string, defaultError string) (reflect.Value, reflect.Value, error) {

	val := reflect.ValueOf(structValue)
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return reflect.Value{}, reflect.Value{}, fmt.Errorf("Bad value, expected struct value, got ")
	}

	var keyVal1 reflect.Value
	var keyVal2 reflect.Value

	for i := 0; i < val.NumField(); i++ {
		fv := val.Field(i)
		ft := typ.Field(i)
		vkey := ft.Tag.Get("vkey")

		if key1 == ft.Name || key1 == vkey {
			keyVal1 = fv
		} else if key2 == ft.Name || key2 == vkey {
			keyVal2 = fv
		}
	}

	if !keyVal1.IsValid() || !keyVal2.IsValid() {
		return reflect.Value{}, reflect.Value{}, errors.New("Unable comparing values, both value need to be provided")
	}

	if keyVal1.Type() != keyVal2.Type() {
		return reflect.Value{}, reflect.Value{}, errors.New("Unable comparing values, both value should have the same type")
	}

	return keyVal1, keyVal2, nil
}

func (v Validation) Email(value interface{}, key, defaultError string) error {
	if IsEmpty(value) {
		return nil
	}
	return v.Match(value, key, v.EmailFormat, defaultError)
}

func (v Validation) Phone(value interface{}, key, defaultError string) error {
	if IsEmpty(value) {
		return nil
	}
	return v.Match(value, key, v.PhoneFormat, defaultError)
}

func (v Validation) Date(value interface{}, key, format, layout, defaultError string) error {
	if IsEmpty(value) {
		return nil
	}

	dateStr, found := value.(string)
	if !found {
		return fmt.Errorf("%s is expected of type string", key)
	}

	_, err := time.Parse(layout, dateStr)

	if err != nil {
		if defaultError == "" {
			return fmt.Errorf("%s is expected of format %s", key, format)
		} else {
			return errors.New(defaultError)
		}

	} else {
		return nil
	}
}

func (v Validation) Match(value interface{}, key, format, defaultError string) error {
	re := regexp.MustCompile(format)
	svalue, found := value.(string)

	if !found {
		return fmt.Errorf("Invalid type, expected string found %s", reflect.TypeOf(value))
	}

	if !re.MatchString(svalue) {
		if defaultError == "" {
			return fmt.Errorf("%s has invalid format value", key)
		}
		return errors.New(defaultError)
	} else {
		return nil
	}
}
