package prtglaw

import "github.com/rs/zerolog"

// Information needed for starting up the application
type Settings struct {
	Host        string
	Port        uint16
	Endpoint    string
	WorkspaceID string
	PrimaryKey  string
	LogTable    string
	LogLevel	string
	Logger 		zerolog.Logger
}

// Represents the data as it is received from PRTG in application/x-www-form-urlencoded format
type DataPRTG struct {
	DateTime    string `schema:"datetime"`
	Date        string `schema:"date"`
	Time        string `schema:"time"`
	TimeZone    string `schema:"timezone"`
	Since       string `schema:"since"`
	Probe       string `schema:"probe"`
	ProbeID     string `schema:"probeid"`
	Group       string `schema:"group"`
	GroupID     string `schema:"groupid"`
	Device      string `schema:"device"`
	DeviceID    string `schema:"deviceid"`
	Sensor      string `schema:"sensor"`
	SensorID    string `schema:"sensorid"`
	Settings    string `schema:"settings"`
	Priority    string `schema:"priority"`
	Host        string `schema:"host"`
	Message     string `schema:"message"`
	LastMessage string `schema:"lastmessage"`
	Status      string `schema:"status"`
	LastStatus  string `schema:"laststatus"`
	LastValue   string `schema:"lastvalue"`
	Downtime    string `schema:"downtime"`
	Uptime      string `schema:"uptime"`
	ObjectTags  string `schema:"objecttags"`
	ParentTags  string `schema:"parenttags"`
}

// Represents the data as it is sent to Microsoft log analytics workspace
type DataKQL struct {
	DateTime    string   `json:"prtg_datetime"`
	Date        string   `json:"prtg_date"`
	Time        string   `json:"prtg_time"`
	TimeZone    string   `json:"prtg_timezone"`
	Since       string   `json:"prtg_since"`
	Probe       string   `json:"prtg_probe"`
	ProbeID     int      `json:"prtg_probeid"`
	Group       string   `json:"prtg_group"`
	GroupID     int      `json:"prtg_groupid"`
	Device      string   `json:"prtg_device"`
	DeviceID    int      `json:"prtg_deviceid"`
	Sensor      string   `json:"prtg_sensor"`
	SensorID    int      `json:"prtg_sensorid"`
	Settings    string   `json:"prtg_settings"`
	Priority    string   `json:"prtg_priority"`
	Host        string   `json:"prtg_host"`
	Message     string   `json:"prtg_message"`
	LastMessage string   `json:"prtg_lastmessage"`
	Status      string   `json:"prtg_status"`
	LastStatus  string   `json:"prtg_laststatus"`
	LastValue   string   `json:"prtg_lastvalue"`
	Downtime    float32  `json:"prtg_downtime"`
	Uptime      float32  `json:"prtg_uptime"`
	ObjectTags  []string `json:"prtg_objecttags"`
	ParentTags  []string `json:"prtg_parenttags"`
}
