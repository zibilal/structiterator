package validator

import (
	"errors"
	"reflect"
	"strings"
)

func IsEmpty(val interface{}) bool {
	return val == nil || reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface())
}

type ValidStruct struct {
	validation Validation
}

func NewValidStruct() *ValidStruct {
	v := Validation{}
	v.PhoneFormat = `^([62]|[0])[0-9]$`
	v.EmailFormat = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	return &ValidStruct{v}
}

func (s *ValidStruct) Valid(input interface{}) []error {
	v := reflect.Indirect(reflect.ValueOf(input))
	t := v.Type()

	var resultError []error

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			fv := v.Field(i)
			ft := t.Field(i)

			if ft.Type.Kind() == reflect.Struct {
				s.Valid(fv)
			} else {
				// process tags
				dtags := ft.Tag.Get("valid")
				vkey := ft.Tag.Get("vkey")
				if dtags != "" {
					dataTags := []*dataTag{}
					dataTags = fetchDataTag(dtags, -1, dataTags)

					for _, dtag := range dataTags {

						valOf := reflect.ValueOf(s.validation)

						if dtag.funcVal != "" {
							val := valOf.MethodByName(dtag.funcVal)
							if val != (reflect.Value{}) {
								if val.IsValid() && val.Type().String() == "func(interface {}, string, string) error" {
									var keyName string
									if vkey != "" {
										keyName = vkey
									} else {
										keyName = ft.Name
									}

									reVal := val.Call([]reflect.Value{
										fv,
										reflect.ValueOf(keyName),
										reflect.ValueOf(dtag.errorMessage),
									})

									if err := processOutput(reVal[0]); err != nil {
										resultError = append(resultError, err)
									}
								} else if val.IsValid() && val.Type().String() == "func(interface {}, string, string, string) error" {
									var keyName string
									if vkey != "" {
										keyName = vkey
									} else {
										keyName = ft.Name
									}
									reVal := val.Call([]reflect.Value{
										fv,
										reflect.ValueOf(keyName),
										reflect.ValueOf(dtag.keyCompare1),
										reflect.ValueOf(dtag.keyCompare2),
										reflect.ValueOf(dtag.errorMessage),
									})

									if err := processOutput(reVal[0]); err != nil {
										resultError = append(resultError, err)
									}
								}
							}

						}
					}
				}
			}
		}

		if len(resultError) > 0 {
			return resultError
		} else {
			return nil
		}

	default:
		return []error{errors.New("Valid only accept input type struct")}
	}
}

func processOutput(reVal reflect.Value) error {
	if !IsEmpty(reVal.Interface()) {
		reErr := reVal.Interface().(error)
		return reErr
	}

	return nil
}

// funcval, errorMessage, key
// funcval, format, errorMessage, key
// funcval, errMessage, key1, key2

type dataTag struct {
	funcVal      string
	errorMessage string
	format       string
	keyCompare1  string
	keyCompare2  string
}

// fetchDataTag idx must always starts from -1
func fetchDataTag(input string, idx int, dataTags []*dataTag) []*dataTag {
	if input == "" {
		return dataTags
	}
	tagsSplits := strings.Split(input, ";")
	if len(tagsSplits) > 1 {
		for i := 0; i < len(tagsSplits); i++ {
			dataTags = append(dataTags, &dataTag{})
		}
		for i := 0; i < len(tagsSplits); i++ {
			fetchDataTag(tagsSplits[i], i, dataTags)
		}
	} else {
		atagSplits := strings.Split(input, ",")
		if len(atagSplits) > 1 {
			if len(dataTags) == 0 {
				dataTags = append(dataTags, &dataTag{})
				idx = 0
			}
			for i := 0; i < len(atagSplits); i++ {
				fetchDataTag(atagSplits[i], idx, dataTags)
			}
		} else {
			splits := strings.Split(input, ":")
			if len(dataTags) == 0 {
				dataTags = append(dataTags, &dataTag{})
				idx = 0
			}
			if len(splits) > 1 {
				itag := dataTags[idx]
				if itag == nil {
					itag = &dataTag{}
					dataTags[idx] = itag
				}
				switch splits[0] {
				case "funcVal":
					itag.funcVal = splits[1]
				case "errorMessage":
					itag.errorMessage = splits[1]
				case "format":
					itag.keyCompare2 = splits[1]
				case "keyCompare1":
					itag.keyCompare1 = splits[1]
				case "keyCompare2":
					itag.keyCompare2 = splits[1]
				}
				fetchDataTag("", idx, dataTags)
			}
		}
	}

	return dataTags

}
