package validator

import (
	"errors"
	"reflect"
	"strings"
)

func IsEmpty(val interface{}) bool {
	empty := val == nil || reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface())

	if !empty {
		// check again for empty string
		if reflect.ValueOf(val).Kind() == reflect.String {
			if strings.TrimSpace(val.(string)) == "" {
				empty = true
			}
		}
	}

	return empty
}

type ValidStruct struct {
	mapper          *ValidationMapper
	PhoneFormat     string
	EmailFormat     string
	DateLayout      string
	DateFormat      string
	ErrorMessageMap map[string]string
}

func NewValidStruct(mapper *ValidationMapper) *ValidStruct {
	v := ValidStruct{}
	v.PhoneFormat = PhoneFormat
	v.EmailFormat = EmailFormat
	v.DateLayout = DateLayout
	v.DateFormat = DateFormat
	v.ErrorMessageMap = make(map[string]string)
	v.mapper = mapper
	v.setupDefaultMapper()
	return &v
}

func NewValidStructWithMap(mapper *ValidationMapper, errorMap map[string]string) *ValidStruct {
	v := NewValidStruct(mapper)
	v.ErrorMessageMap = errorMap
	return v
}

func (s *ValidStruct) setupDefaultMapper() error {
	ival := reflect.ValueOf(Validation{})
	for i := 0; i < ival.NumMethod(); i++ {
		method := ival.Method(i)
		mtype := method.Type()
		if err := s.mapper.AddFunc(mtype.Name(), method.Interface()); err != nil {
			return err
		}
	}

	return nil
}

func (s *ValidStruct) RegisterValidator(name string, f interface{}) error {
	return s.mapper.AddFunc(name, f)
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
				jsplit := strings.Split(ft.Tag.Get("json"), ",")
				vkey := ""
				if len(jsplit) > 0 {
					vkey = strings.TrimSpace(jsplit[0])
				}

				if dtags != "" {
					dataTags := []*dataTag{}
					dataTags = fetchDataTag(dtags, -1, dataTags)

					for _, dtag := range dataTags {
						if len(s.ErrorMessageMap) > 0 && dtag.errorMessage == "" {
							tmp, found := s.ErrorMessageMap[dtag.funcVal]
							if found {
								dtag.errorMessage = tmp
							}
						}

						if dtag.funcVal != "" {
							ival, err := s.mapper.GetFunc(dtag.funcVal)
							if err != nil {
								resultError = append(resultError)
								return resultError
							}
							val := reflect.ValueOf(ival)
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
									k1, k2 := "", ""
									var theValue reflect.Value
									if dtag.compareKey != "" && dtag.compareValue != "" {
										k1, k2 = dtag.compareKey, dtag.compareValue
										theValue = v
									} else if dtag.compareKey != "" && dtag.compareValue == "" {
										var keyName string
										if vkey != "" {
											keyName = vkey
										} else {
											keyName = ft.Name
										}
										k1, k2 = keyName, dtag.compareKey
										theValue = v
									} else if dtag.acceptedValues != "" {
										theValue = fv
										if vkey != "" {
											k1 = vkey
										} else {
											k1 = ft.Name
										}
										k2 = dtag.acceptedValues
									} else if dtag.format != "" {
										theValue = fv
										if vkey != "" {
											k1 = vkey
										} else {
											k1 = ft.Name
										}
										k2 = dtag.format
									}

									if k1 != "" && k2 != "" {
										reVal := val.Call([]reflect.Value{
											theValue,
											reflect.ValueOf(k1),
											reflect.ValueOf(k2),
											reflect.ValueOf(dtag.errorMessage),
										})

										if err := processOutput(reVal[0]); err != nil {
											resultError = append(resultError, err)
										}
									}

								} else if val.IsValid() && val.Type().String() == "func(interface {}, string, interface {}, string, string, string) error" {
									k1, k2 := "", ""
									if dtag.compareKey != "" && dtag.compareValue != "" {
										k1, k2 = dtag.compareKey, dtag.compareValue
									}

									if k1 != "" && k2 != "" {
										var keyName string
										if vkey != "" {
											keyName = vkey
										} else {
											keyName = ft.Name
										}

										reVal := val.Call([]reflect.Value{
											v,
											reflect.ValueOf(keyName),
											fv,
											reflect.ValueOf(k1),
											reflect.ValueOf(k2),
											reflect.ValueOf(dtag.errorMessage),
										})

										if err := processOutput(reVal[0]); err != nil {
											resultError = append(resultError, err)
										}
									}
								} else if val.IsValid() && val.Type().String() == "func(interface {}, string, string, string, string) error" {
									k1, k2 := "", ""
									if dtag.format != "" && dtag.dateLayout != "" {
										k1, k2 = dtag.format, dtag.dateLayout
									} else {
										k1, k2 = s.DateFormat, s.DateLayout
									}

									if k1 != "" && k2 != "" {
										var keyName string
										if vkey != "" {
											keyName = vkey
										} else {
											keyName = ft.Name
										}

										reVal := val.Call([]reflect.Value{
											fv,
											reflect.ValueOf(keyName),
											reflect.ValueOf(k1),
											reflect.ValueOf(k2),
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
		}

		if len(resultError) > 0 {
			return resultError
		} else {
			return nil
		}

	default:
		return []error{errors.New("valid only accept input type struct")}
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
	funcVal        string
	errorMessage   string
	format         string
	compareKey     string
	compareValue   string
	dateLayout     string
	acceptedValues string
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
					itag.format = splits[1]
				case "compareValue":
					itag.compareValue = splits[1]
				case "compareKey":
					itag.compareKey = splits[1]
				case "dateLayout":
					itag.dateLayout = splits[1]
				case "values":
					itag.acceptedValues = splits[1]
				}
				fetchDataTag("", idx, dataTags)
			}
		}
	}

	return dataTags

}
