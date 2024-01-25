package helper

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var Errors map[string]map[string]string

var mu sync.Mutex

func SetError(b *testing.B, orm, method, err string) {
	b.Helper()

	mu.Lock()
	Errors[orm][method] = err
	mu.Unlock()
	b.Fail()
}

func GetError(orm, method string) string {
	return Errors[orm][method]
}

func getFuncName(function interface{}) string {
	name := strings.Split(runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name(), ".")
	straightName := strings.Split(name[len(name)-1], "-")[0]

	return straightName
}
