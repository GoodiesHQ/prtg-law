package prtglaw

import (
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	DATETIME_LAYOUT_PRTG = "1/2/2006 3:04:05 PM"
	DATETIME_LAYOUT_ISO  = "2006-01-02T15:04:05"
	DATE_LAYOUT_PRTG     = "1/2/2006"
	DATE_LAYOUT_ISO      = "2006-01-02"
	TIME_LAYOUT_PRTG     = "3:04:05 PM"
	TIME_LAYOUT_ISO      = "15:04:05"
)

// convert from PRTG datetime to ISO format without timezone information
func datetimePRTGtoISO(dt string, fallback string) string {
	val, err := time.Parse(DATETIME_LAYOUT_PRTG, dt)
	if err != nil {
		return fallback
	}
	return val.Format(DATETIME_LAYOUT_ISO)
}

// convert from PRTG datetime to ISO format without timezone information
func datePRTGtoISO(d string, fallback string) string {
	val, err := time.Parse(DATE_LAYOUT_PRTG, d)
	if err != nil {
		return fallback
	}
	return val.Format(DATE_LAYOUT_ISO)
}

// convert from PRTG datetime to ISO format without timezone information
func timePRTGtoISO(t string, fallback string) string {
	val, err := time.Parse(TIME_LAYOUT_PRTG, t)
	if err != nil {
		return fallback
	}
	return val.Format(TIME_LAYOUT_ISO)
}

func toInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}

func toFloat32(s string) float32 {
	if f, err := strconv.ParseFloat(s, 32); err == nil {
		return float32(f)
	}
	return 0.0
}

func isEmpty(s string) bool {
	return s == ""
}

func PRTGtoKQL(prtg *DataPRTG) (*DataKQL, error) {
	var result = new(DataKQL)

	// process date and time related information
	result.DateTime = datetimePRTGtoISO(prtg.DateTime, prtg.DateTime)
	result.Date = datePRTGtoISO(prtg.Date, prtg.Date)
	result.Time = timePRTGtoISO(prtg.Time, prtg.Time)
	result.TimeZone = prtg.TimeZone
	result.Since = datetimePRTGtoISO(prtg.Since, prtg.Since)

	// process info about the sensor
	result.Probe = prtg.Probe
	result.ProbeID = toInt(prtg.ProbeID)
	result.Group = prtg.Group
	result.GroupID = toInt(prtg.GroupID)
	result.Device = prtg.Device
	result.DeviceID = toInt(prtg.DeviceID)
	result.Sensor = prtg.Sensor
	result.SensorID = toInt(prtg.SensorID)
	result.Priority = prtg.Priority
	result.Host = prtg.Host

	// process info about the event
	result.Message = prtg.Message
	result.LastMessage = prtg.LastMessage
	result.Status = prtg.Status
	result.LastStatus = prtg.Status
	result.LastValue = prtg.LastValue

	if idx := strings.Index(prtg.Downtime, "%"); idx >= 0 {
		result.Downtime = toFloat32(prtg.Downtime[:idx])
	}

	if idx := strings.Index(prtg.Uptime, "%"); idx >= 0 {
		result.Uptime = toFloat32(prtg.Uptime[:idx])
	}

	result.ObjectTags = slices.DeleteFunc(strings.Split(prtg.ObjectTags, " "), isEmpty)
	result.ParentTags = slices.DeleteFunc(strings.Split(prtg.ParentTags, " "), isEmpty)

	return result, nil
}
