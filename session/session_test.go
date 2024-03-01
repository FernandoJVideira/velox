package session

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alexedwards/scs/v2"
)

func TestSession_InitSession(t *testing.T) {
	v := &Session{
		CookieLifetime: "60",
		CookiePersist:  "true",
		CookieName:     "velox",
		CookieDomain:   "localhost",
		SessionType:    "cookie",
	}
	var sm *scs.SessionManager
	ses := v.InitSession()

	var sessKind reflect.Kind
	var sessType reflect.Type

	rv := reflect.ValueOf(ses)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		fmt.Println("For loop: ", rv.Kind(), rv.Type(), rv)
		sessKind = rv.Kind()
		sessType = rv.Type()

		rv = rv.Elem()
	}

	if !rv.IsValid() {
		t.Error("invalid type or kind; kind:", rv.Kind(), "type:", rv.Type())
	}

	if sessKind != reflect.ValueOf(sm).Kind() {
		t.Error("wrong kind testing cookie session; expected:", reflect.ValueOf(sm).Kind(), "got:", sessKind)
	}
	if sessType != reflect.ValueOf(sm).Type() {
		t.Error("wrong type testing cookie session; expected:", reflect.ValueOf(sm).Kind(), "got:", sessKind)
	}

}
