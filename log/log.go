package log

import (
	"encoding/json"
	"fmt"
	llog "log"
	"net/http"
	"strings"
)

var ProjectID string

func init() {
	llog.SetFlags(0)
}

type Entry struct {
	Severity string `json:"severity,omitempty"`
	Message  string `json:"message"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`

	Component string `json:"component"`
}

func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	if e.Component == "" {
		e.Component = "default"
	}
	out, err := json.Marshal(e)
	if err != nil {
		llog.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

func Debugf(r *http.Request, m string, args ...interface{}) {
	Log(r, "DEBUG", m, args...)
}

func Infof(r *http.Request, m string, args ...interface{}) {
	Log(r, "INFO", m, args...)
}

func Warningf(r *http.Request, m string, args ...interface{}) {
	Log(r, "WARNING", m, args...)
}

func Criticalf(r *http.Request, m string, args ...interface{}) {
	Log(r, "CRITICAL", m, args...)
}

func Errorf(r *http.Request, m string, args ...interface{}) {
	Log(r, "ERROR", m, args...)
}

func Fatalf(w http.ResponseWriter, r *http.Request, m string, args ...interface{}) {
	Warningf(r, m, args...)
	http.Error(w, fmt.Sprintf(m, args...), http.StatusInternalServerError)
}

func Log(r *http.Request, severity string, m string, args ...interface{}) {
	// Detect the log trace from request header, if available.
	var trace string
	if r != nil {
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			trace = fmt.Sprintf("projects/%s/traces/%s", ProjectID, traceParts[0])
		}
	}

	llog.Println(Entry{
		Severity: severity,
		Message:  fmt.Sprintf(m, args...),
		Trace:    trace,
	})
}
