package validator

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ValidationMapper struct {
	funcMap            map[string]interface{}
	acceptedSignatures []string
	sync.Mutex
}

func NewValidationMapper() *ValidationMapper {
	vMapper := new(ValidationMapper)
	vMapper.funcMap = make(map[string]interface{})
	vMapper.acceptedSignatures = []string{
		"func(interface {}, string, string) error",
		"func(interface {}, string, string, string) error",
		"func(interface {}, string, interface {}, string, string, string) error",
		"func(interface {}, string, string, string, string) error",
	}

	return vMapper
}

func (v *ValidationMapper) AddFunc(name string, f interface{}) error {
	fValue := reflect.ValueOf(f)
	if fValue.Kind() != reflect.Func {
		return errors.New("please provide a function typed argument")
	}

	notFound := true
	for _, s := range v.acceptedSignatures {
		if fValue.Type().String() == s {
			notFound = false
			break
		}
	}

	if notFound {
		return errors.New("function accepted is not accepted")
	}

	v.Lock()
	v.funcMap[name] = f
	v.Unlock()

	return nil
}

func (v *ValidationMapper) GetFunc(name string) (interface{}, error) {
	var (
		result interface{}
		found  bool
	)
	v.Lock()
	result, found = v.funcMap[name]
	v.Unlock()

	if found {
		return result, nil
	} else {
		return nil, fmt.Errorf("func name %s is not found", name)
	}
}

const (
	PhoneFormat = `^([62]|[0])[0-9]+$`
	EmailFormat = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	DateLayout  = `01/02/2006`
	DateFormat  = `mm/dd/yyyy`
)

type Validation struct {
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

func (v Validation) CondRequired(structValue interface{}, key string, zeValue interface{}, keyCompare, valueCompare string, defaultError string) error {

	val := reflect.ValueOf(structValue)
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("bad value, expected struct value, got ")
	}

	if keyCompare == "" && valueCompare == "" {
		return errors.New("bad state, keyCompare and valueCompare is expected to have a string value")
	}

	for i := 0; i < val.NumField(); i++ {
		fv := val.Field(i)
		ft := typ.Field(i)

		vkey := ft.Tag.Get("json")
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
		return fmt.Errorf("expected type string got %s", val1.Type())
	}

	if ival1 == "" {
		return nil
	}

	ival2, found := val2.Interface().(string)
	if !found {
		return fmt.Errorf("expected type string got %s", val2.Type())
	}

	var dateLayout = DateLayout

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
			return fmt.Errorf("invalid %s should be after %s", key1, key2)
		}

		return errors.New(defaultError)
	}

	return nil
}

func (v Validation) after(structValue interface{}, key1, key2 string, defaultError string) (reflect.Value, reflect.Value, error) {

	val := reflect.ValueOf(structValue)
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return reflect.Value{}, reflect.Value{}, fmt.Errorf("bad value, expected struct value, got ")
	}

	var keyVal1 reflect.Value
	var keyVal2 reflect.Value

	for i := 0; i < val.NumField(); i++ {
		fv := val.Field(i)
		ft := typ.Field(i)
		vkey := ft.Tag.Get("json")

		if key1 == ft.Name || key1 == vkey {
			keyVal1 = fv
		} else if key2 == ft.Name || key2 == vkey {
			keyVal2 = fv
		}
	}

	if !keyVal2.IsValid() {
		return reflect.Value{}, reflect.Value{}, errors.New("unable comparing values, both value need to be provided")
	}

	if keyVal1.Type() != keyVal2.Type() {
		return reflect.Value{}, reflect.Value{}, errors.New("unable comparing values, both value should have the same type")
	}

	return keyVal1, keyVal2, nil
}

func (v Validation) Email(value interface{}, key, defaultError string) error {
	if IsEmpty(value) {
		return nil
	}
	return v.Match(value, key, EmailFormat, defaultError)
}

func (v Validation) Url(value interface{}, key, defaultError string) error {
	if IsEmpty(value) {
		return nil
	}

	str, ok := value.(string)

	if ok {
		_, err := url.ParseRequestURI(str)

		if err != nil {
			if defaultError == "" {
				defaultError = err.Error()
			}
			return errors.New(defaultError)
		}
	}

	return nil
}

func (v Validation) Phone(value interface{}, key, defaultError string) error {
	if IsEmpty(value) {
		return nil
	}
	return v.Match(value, key, PhoneFormat, defaultError)
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

	if IsEmpty(value) {
		return nil
	}

	re, err := regexp.Compile(format)
	if err != nil {
		return fmt.Errorf("invalid regular expression %s: %s", format, err.Error())
	}
	svalue, found := value.(string)

	if !found {
		return fmt.Errorf("invalid type, expected string found %s", reflect.TypeOf(value))
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

func (v Validation) AcceptedValues(value interface{}, key, theValues, defaultError string) error {

	if IsEmpty(value) {
		return nil
	}

	var vs []string
	var isRange bool
	if strings.Contains(theValues, "<->") {
		vs = strings.Split(theValues, "<->")
		isRange = true
	} else if strings.Contains(theValues, "|") {
		vs = strings.Split(theValues, "|")
		isRange = false
	}

	if isRange {
		return checkInRange(value, vs[0], vs[1], defaultError)
	} else {
		return checkInValues(value, vs, theValues, defaultError)
	}

	return nil
}

func checkInValues(val interface{}, vals []string, acceptedValues, errorMessage string) error {
	for _, v := range vals {
		switch val.(type) {
		case int:
			ival := val.(int)
			iv, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			if ival == int(iv) {
				return nil
			}
		case int64:
			ival := val.(int64)
			iv, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			if ival == int64(iv) {
				return nil
			}
		case uint:
			ival := val.(uint)
			iv, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			if ival == uint(iv) {
				return nil
			}
		case uint64:
			ival := val.(uint64)
			iv, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			if ival == iv {
				return nil
			}
		case float32:
			ival := val.(float32)
			iv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			if ival == float32(iv) {
				return nil
			}
		case float64:
			ival := val.(float64)
			iv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			if ival == iv {
				return nil
			}
		case string:
			ival := val.(string)
			if ival == v {
				return nil
			}
		}
	}

	if errorMessage == "" {
		return fmt.Errorf("wrong value %v, accepted values %s", val, acceptedValues)
	} else {
		return errors.New(errorMessage)
	}
}

func checkInRange(val interface{}, val1, val2 string, errorMessage string) error {

	switch val.(type) {
	case int:
		ival := val.(int)
		ival1, err := strconv.ParseInt(val1, 10, 64)
		if err != nil {
			return err
		}
		ival2, err := strconv.ParseInt(val2, 10, 64)
		if err != nil {
			return err
		}
		if !(ival >= int(ival1) && ival <= int(ival2)) {
			if errorMessage == "" {
				return fmt.Errorf("%d is outside of range %d - %d", ival, ival1, ival2)
			} else {
				return errors.New(errorMessage)
			}
		}
	case int64:
		ival := val.(int64)
		ival1, err := strconv.ParseInt(val1, 10, 64)
		if err != nil {
			return err
		}
		ival2, err := strconv.ParseInt(val2, 10, 64)
		if err != nil {
			return err
		}
		if !(ival >= ival1 && ival <= ival2) {
			if errorMessage == "" {
				return fmt.Errorf("%d is outside of range %d - %d", ival, ival1, ival2)
			} else {
				return errors.New(errorMessage)
			}
		}
	case uint:
		ival := val.(uint)
		ival1, err := strconv.ParseUint(val1, 10, 64)
		if err != nil {
			return err
		}
		ival2, err := strconv.ParseUint(val2, 10, 64)
		if err != nil {
			return err
		}
		if !(ival >= uint(ival1) && ival <= uint(ival2)) {
			if errorMessage == "" {
				return fmt.Errorf("%d is outside of range %d - %d", ival, ival1, ival2)
			} else {
				return errors.New(errorMessage)
			}
		}
	case uint64:
		ival := val.(uint64)
		ival1, err := strconv.ParseUint(val1, 10, 64)
		if err != nil {
			return err
		}
		ival2, err := strconv.ParseUint(val2, 10, 64)
		if err != nil {
			return err
		}
		if !(ival >= ival1 && ival <= ival2) {
			if errorMessage == "" {
				return fmt.Errorf("%d is outside of range %d - %d", ival, ival1, ival2)
			} else {
				return errors.New(errorMessage)
			}
		}
	default:
		return errors.New("for check in range only accept int|int64|uint|uint64")
	}

	return nil
}
