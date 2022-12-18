/*
 * Copyright (c) 2022 Huawei Device Co., Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
