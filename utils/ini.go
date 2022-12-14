package utils

import (
	"github.com/Unknwon/goconfig"
	"github.com/sirupsen/logrus"
	"reflect"
)

// ParseFromConfigFile parse ini file and set values by the tag of fields.
// 'p' must be a pointer to the given structure, otherwise will panic.
// Only process its string fields and its sub structs.
func ParseFromConfigFile(section string, p any) {
	conf, err := goconfig.LoadConfigFile("fotff.ini")
	if err != nil {
		logrus.Warnf("load config file err: %v", err)
	}
	rv := reflect.ValueOf(p)
	rt := reflect.TypeOf(p)
	for i := 0; i < rv.Elem().NumField(); i++ {
		switch rt.Elem().Field(i).Type.Kind() {
		case reflect.String:
			key := rt.Elem().Field(i).Tag.Get("key")
			if key == "" {
				continue
			}
			var v string
			if conf != nil {
				v, err = conf.GetValue(section, key)
			}
			if conf == nil || err != nil {
				v = rt.Elem().Field(i).Tag.Get("default")
			}
			rv.Elem().Field(i).SetString(v)
		case reflect.Struct:
			ParseFromConfigFile(section, rv.Elem().Field(i).Addr().Interface())
		}
	}
}
