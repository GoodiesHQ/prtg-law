package main

import (
	prtglaw "github.com/goodieshq/prtg-law/pkg"
	"github.com/joho/godotenv"
)

// environment variables
const (
	EV_WORKSPACE_ID  = "WORKSPACE_ID"  // log analytics workspace ID
	EV_PRIMARY_KEY   = "PRIMARY_KEY"   // log analytics workspace primary key
	EV_TABLE_NAME    = "TABLE_NAME"    // the log analytics workspace table name to use (Azure adds a _CL suffix)
	EV_LISTEN_HOST   = "LISTEN_HOST"   // host addr on which to listen
	EV_LISTEN_PORT   = "LISTEN_PORT"   // host port on which to listen
	EV_HTTP_ENDPOINT = "HTTP_ENDPOINT" // the http endpoint to serve the webhook handler on
)

// defaults
var (
	DEFAULT_TABLE_NAME    string = "PRTG"
	DEFAULT_LISTEN_HOST   string = "0.0.0.0"
	DEFAULT_HTTP_ENDPOINT string = "/prtg"
	DEFAULT_LISTEN_PORT   uint16 = 8888
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
}

type Example struct {
	Name    string   `json:"name"`
	Percent float32  `json:"percent"`
	Array   []string `json:"array"`
}

func main() {
	var settings prtglaw.Settings
	env(&settings)
	prtglaw.Serve(&settings)
}
