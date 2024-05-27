package prtglaw

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

const (
	LAW_HTTP_METHOD  = "POST"
	LAW_CONTENT_TYPE = "application/json"
	LAW_API_VERSION  = "2016-04-01"
	LAW_ENDPOINT     = "/api/logs"
)

func sign(primaryKey, rfc1123 string, contentLength int) (string, error) {
	// craft the payload as per https://learn.microsoft.com/en-us/rest/api/loganalytics/create-request#constructing-the-signature-string
	payload := []byte(fmt.Sprintf("%s\n%d\n%s\nx-ms-date:%s\n%s", LAW_HTTP_METHOD, contentLength, LAW_CONTENT_TYPE, rfc1123, LAW_ENDPOINT))

	// decode the primary key from base64
	key, err := base64.StdEncoding.DecodeString(primaryKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode primary key: %w", err)
	}

	// calculate the HMAC-SHA256 digest of the crafted payload
	h := hmac.New(sha256.New, key)
	n, err := h.Write(payload)
	if err != nil {
		return "", fmt.Errorf("failed to digest the payload: %w", err)
	}
	if n != len(payload) {
		return "", fmt.Errorf("failed to digest the full payload")
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}


func PopulateLAW(settings *Settings, payload interface{}) error {
	// calculate the proper RFC1123 date with GMT instead of UTC
	rfc1123 := time.Now().UTC().Format(time.RFC1123)
	rfc1123 = rfc1123[:len(rfc1123)-3] + "GMT"

	// convert the payload into JSON bytes
	body, err := Dumpb(payload, settings)
	if err != nil {
		return err
	}

	// calculate the signature for the Authorization header
	signature, err := sign(settings.PrimaryKey, rfc1123, len(body))
	if err != nil {
		return err
	}

	// craft the URL
	url := fmt.Sprintf("https://%s.ods.opinsights.azure.com/%s?api-version=%s", settings.WorkspaceID, LAW_ENDPOINT, LAW_API_VERSION)
	req, err := http.NewRequest(LAW_HTTP_METHOD, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// set the request headers
	req.Header.Set("Content-Type", LAW_CONTENT_TYPE)
	req.Header.Set("Authorization", fmt.Sprintf("SharedKey %s:%s", settings.WorkspaceID, signature))
	req.Header.Set("x-ms-date", rfc1123)
	req.Header.Set("Log-Type", settings.LogTable)

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	return fmt.Errorf("request failed: %s", resp.Status)
}
