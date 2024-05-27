package main

import (
	"os"
	"strings"
	"time"

	prtglaw "github.com/goodieshq/prtg-law/pkg"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// environment variables
const (
	EV_WORKSPACE_ID  = "WORKSPACE_ID"  // log analytics workspace ID
	EV_PRIMARY_KEY   = "PRIMARY_KEY"   // log analytics workspace primary key
	EV_TABLE_NAME    = "TABLE_NAME"    // the log analytics workspace table name to use (Azure adds a _CL suffix)
	EV_LISTEN_HOST   = "LISTEN_HOST"   // host addr on which to listen
	EV_LISTEN_PORT   = "LISTEN_PORT"   // host port on which to listen
	EV_HTTP_ENDPOINT = "HTTP_ENDPOINT" // the http endpoint to serve the webhook handler on
	EV_LOG_LEVEL     = "LOG_LEVEL"     // the level of logging (DEBUG, INFO, WARNING, ERROR)
)

// defaults
var (
	DEFAULT_TABLE_NAME    string = "PRTG"
	DEFAULT_LISTEN_HOST   string = "0.0.0.0"
	DEFAULT_HTTP_ENDPOINT string = "/prtg"
	DEFAULT_LISTEN_PORT   uint16 = 8888
	DEFAULT_LOG_LEVEL     string = "INFO"
)

func env(settings *prtglaw.Settings) {
	// attempt to load environment variables from a file, but continue on failue
	godotenv.Load()

	// parse the required parameters, exits if not found
	EnvLookupMust(&settings.WorkspaceID, EV_WORKSPACE_ID)
	EnvLookupMust(&settings.PrimaryKey, EV_PRIMARY_KEY)

	// parse the optional parameters or use the defaults
	EnvLookup(&settings.LogTable, EV_TABLE_NAME, &DEFAULT_TABLE_NAME)
	EnvLookup(&settings.Host, EV_LISTEN_HOST, &DEFAULT_LISTEN_HOST)
	EnvLookup(&settings.Port, EV_LISTEN_PORT, &DEFAULT_LISTEN_PORT)
	EnvLookup(&settings.Endpoint, EV_HTTP_ENDPOINT, &DEFAULT_HTTP_ENDPOINT)
	EnvLookup(&settings.LogLevel, EV_LOG_LEVEL, &DEFAULT_LOG_LEVEL)
}

type Example struct {
	Name    string   `json:"name"`
	Percent float32  `json:"percent"`
	Array   []string `json:"array"`
}

func main() {
	var settings prtglaw.Settings
	env(&settings)

	settings.LogLevel = strings.ToUpper(settings.LogLevel)
	level := func(levelName string) zerolog.Level {
		switch levelName {
		case "DEBUG":
			return zerolog.DebugLevel
		case "INFO":
			return zerolog.InfoLevel
		case "WARN":
			return zerolog.WarnLevel
		case "ERROR":
			return zerolog.WarnLevel
		default:
			return zerolog.InfoLevel
		}
	}(settings.LogLevel)

	settings.Logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(level).With().Timestamp().Logger()

	prtglaw.Serve(&settings)
}
