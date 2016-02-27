package logging

import (
	"os"

	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("vasuki")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfile}#%{shortfunc} ▶ %{level:.10s} %{id:04d}%{color:reset} %{message}`,
)

func init() {
	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	formattedStdout := logging.NewBackendFormatter(stdout, format)
	logging.SetBackend(formattedStdout)
	logging.SetLevel(logging.INFO, "")
}

func EnableDebug(enable bool) {
	if enable {
		logging.SetLevel(logging.DEBUG, "")
	}
}
