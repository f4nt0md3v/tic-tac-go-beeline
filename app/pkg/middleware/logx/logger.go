package logx

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/middleware"
	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/netx"
)

// LoggerConfig defines the config for Logger middleware.
type LogConfig struct {
	Formatter LogFormatter
}

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	Request    *http.Request
	TimeStamp  time.Time
	StatusCode int
	Latency    time.Duration
	ClientIP   net.IP
	Method     string
	Path       string
}

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFmt = func(param LogFormatterParams) string {
	var (
		methodColor = param.MethodColor()
		resetColor  = param.ResetColor()
	)

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("%v | %s %s %s | [%s] | %#v | %v |",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		methodColor, param.Method, resetColor,
		param.ClientIP,
		param.Path,
		param.Latency,
	)
}

// Logger instances a Logger middleware that will write the logs for all incoming requests
// By default to os.Stdout.
func Logger() middleware.Middleware {
	formatter := defaultLogFmt
	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			path := r.URL.Path
			raw := r.URL.RawQuery

			param := LogFormatterParams{
				Request: r,
			}

			defer func() { log.Println(formatter(param)) }()

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = netx.GetClientIPFromRequest(r)
			param.Method = r.Method

			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path

			f(w, r)
		}
	}
}
