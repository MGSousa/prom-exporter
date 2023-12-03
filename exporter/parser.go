package exporter

import (
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
)

// parse extracted raw data and normalize it
func (e *Exporter) parse(v any, keys ...string) {
	var (
		val interface{}
	)

	switch x := getType(v); x {
	case reflect.Map:
		for k, nv := range v.(map[string]any) {
			e.parse(nv, append(keys, k)...)
		}
		return
	case
		reflect.Float64,
		reflect.Float32,
		reflect.Int64,
		reflect.Int32:
		val = v
	case reflect.String:
		val = fmt.Sprintf("%s::untyped", v)
	default:
		log.Error("parsed type is unknown", x)
	}

	e.data = append(e.data,
		fmt.Sprintf("%s::%s::%v", join(keys[:len(keys)-1], "_"), keys[len(keys)-1], val))
}

func getType(t any) reflect.Kind {
	return reflect.TypeOf(t).Kind()
}
