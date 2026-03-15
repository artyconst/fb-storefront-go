package storefront

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	sf "github.com/artyconst/fb-storefront-go/internal/types"
)

// StdLogger implements Logger using the standard log package.
type StdLogger struct {
	logger *log.Logger
	level  sf.LogLevel
}

// NewStdLogger creates a new standard logger with the given output and level.
func NewStdLogger(output io.Writer, level sf.LogLevel) *StdLogger {
	return &StdLogger{
		logger: log.New(output, "", 0),
		level:  level,
	}
}

func (sl *StdLogger) shouldLog(level sf.LogLevel) bool {
	return int(level) >= int(sl.level)
}

func (sl *StdLogger) print(level sf.LogLevel, format string, v ...interface{}) {
	if !sl.shouldLog(level) {
		return
	}
	prefix := fmt.Sprintf("[%s]", level.String())
	sl.logger.Output(2, fmt.Sprintf("%s %s", prefix, fmt.Sprintf(format, v...)))
}

func (sl *StdLogger) Error(v ...interface{}) {
	if sl.shouldLog(sf.LevelError) {
		sl.logger.Println(v...)
	}
}

func (sl *StdLogger) Errorf(format string, v ...interface{}) {
	sl.print(sf.LevelError, format, v...)
}

func (sl *StdLogger) Warn(v ...interface{}) {
	if sl.shouldLog(sf.LevelWarn) {
		prefix := fmt.Sprintf("[%s]", sf.LevelWarn.String())
		sl.logger.Output(2, fmt.Sprintf("%s %v", prefix, v))
	}
}

func (sl *StdLogger) Warnf(format string, v ...interface{}) {
	sl.print(sf.LevelWarn, format, v...)
}

func (sl *StdLogger) Info(v ...interface{}) {
	if sl.shouldLog(sf.LevelInfo) {
		prefix := fmt.Sprintf("[%s]", sf.LevelInfo.String())
		sl.logger.Output(2, fmt.Sprintf("%s %v", prefix, v))
	}
}

func (sl *StdLogger) Infof(format string, v ...interface{}) {
	sl.print(sf.LevelInfo, format, v...)
}

func (sl *StdLogger) Debug(v ...interface{}) {
	if sl.shouldLog(sf.LevelDebug) {
		prefix := fmt.Sprintf("[%s]", sf.LevelDebug.String())
		sl.logger.Output(2, fmt.Sprintf("%s %v", prefix, v))
	}
}

func (sl *StdLogger) Debugf(format string, v ...interface{}) {
	sl.print(sf.LevelDebug, format, v...)
}

// RawResponseLogger handles logging of raw HTTP request/response bodies.
type RawResponseLogger struct {
	logger            *StdLogger
	annotate          bool
	enableRequestBody bool
	requestIndent     string
	responseIndent    string
	borderLine        string
}

// NewRawResponseLogger creates a new raw response logger.
func NewRawResponseLogger(logger *StdLogger, config *sf.LoggingConfig) *RawResponseLogger {
	return &RawResponseLogger{
		logger:            logger,
		annotate:          config.AnnotateRawResponses,
		enableRequestBody: config.EnableRequestBody,
		requestIndent:     "  ",
		responseIndent:    "  ",
		borderLine:        strings.Repeat("-", 50),
	}
}

// LogRequest logs the raw HTTP request details.
func (r *RawResponseLogger) LogRequest(method, url string, body interface{}) {
	if !r.enableRequestBody && method != http.MethodGet {
		r.logger.Debug("REQUEST: %s %s", method, url)
		return
	}

	r.logger.Debug("%s\n%s REQUEST %s\n%s", r.borderLine, method, url, r.borderLine)
	if body != nil {
		r.logBody("BODY", body, r.requestIndent)
	}
}

// LogResponse logs the raw HTTP response details.
func (r *RawResponseLogger) LogResponse(statusCode int, body []byte) {
	r.logger.Debug("%s\n%s RESPONSE %d\n%s", r.borderLine, statusCode, r.borderLine)
	if len(body) > 0 {
		r.logBody("BODY", string(body), r.responseIndent)
	}
}

// logBody logs a JSON body with proper formatting.
func (r *RawResponseLogger) logBody(label string, data interface{}, indent string) {
	var formatted string

	if s, ok := data.(string); ok {
		formatted = s
	} else if b, ok := data.([]byte); ok {
		formatted = string(b)
	} else {
		jsonBytes, err := json.MarshalIndent(data, indent, "  ")
		if err != nil {
			r.logger.Debug("%s %v", label, data)
			return
		}
		formatted = string(jsonBytes)
	}

	if r.annotate {
		r.logger.Debug("\n%s\n%s\n%s", strings.Repeat("=", 50), formatted, strings.Repeat("=", 50))
	} else {
		lines := strings.Split(formatted, "\n")
		for _, line := range lines {
			if line != "" {
				r.logger.Debug("%s %s", indent, line)
			}
		}
	}
}

// PrettyPrintJSON formats JSON for console output with annotations.
func PrettyPrintJSON(data *sf.LoggingConfig) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", data)
	}

	lines := strings.Split(string(jsonBytes), "\n")
	result := []string{"\n" + strings.Repeat("=", 60)}
	for _, line := range lines {
		result = append(result, "│ "+line)
	}
	result = append(result, "└"+strings.Repeat("─", 58))

	return strings.Join(result, "\n")
}

// PrettyPrintRaw logs raw JSON with console-friendly annotations.
func PrettyPrintRaw(rawJSON string) {
	if len(rawJSON) == 0 {
		fmt.Println("\n[Empty response body]")
		return
	}

	lines := strings.Split(rawJSON, "\n")
	result := []string{"\n" + strings.Repeat("=", 60)}
	for _, line := range lines {
		if len(line) > 58 {
			line = line[:55] + "..."
		}
		result = append(result, "│ "+line)
	}
	result = append(result, "└"+strings.Repeat("─", 58))

	fmt.Println(strings.Join(result, "\n"))
}
