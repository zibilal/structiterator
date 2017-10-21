package validator

import (
	"reflect"
	"errors"
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
				if dtags != "" {
					tags := strings.Split(dtags, ",")
					var (
						funcVal string
						errMessage string
					)

					if len(tags) == 2 {
						funcVal = tags[0]
						errMessage = tags[1]
					} else {
						funcVal = tags[0]
					}

					valOf := reflect.ValueOf(s.validation)

					val := valOf.MethodByName(funcVal)
					if val != (reflect.Value{}) {
						if val.IsValid() && val.Type().String() == "func(interface {}, string, string) error" {
							reVal := val.Call([]reflect.Value{
								fv,
								reflect.ValueOf(ft.Name),
								reflect.ValueOf(errMessage),
							})

							if !IsEmpty(reVal[0].Interface()) {
								reErr := reVal[0].Interface().(error)
								resultError = append(resultError, reErr)
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
		return []error {errors.New("Valid only accept input type struct")}
	}
}
