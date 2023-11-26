package exporter

import (
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Parser struct {
	data []string
}

func (p *Parser) parse(v any, keys ...string) {
	var val interface{}

	switch x := getType(v); x {
	case reflect.Map:
		for k, nv := range v.(map[string]any) {
			p.parse(nv, append(keys, k)...)
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

	metric := strings.Join(keys[:len(keys)-1], "_")

	// TODO: change this
	p.data = append(p.data, fmt.Sprintf("%s::%s::%v", metric, keys[len(keys)-1], val))
}

func getType(t any) reflect.Kind {
	return reflect.TypeOf(t).Kind()
}

func split(s, delim string) []string {
	return strings.Split(s, delim)
}

func join(s []string, delim string) string {
	return strings.Join(s, delim)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
