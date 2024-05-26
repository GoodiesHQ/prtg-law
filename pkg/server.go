package prtglaw

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

var decoder = schema.NewDecoder()

// PRTG webhook handler

func makePRTGHandler(settings *Settings) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("Webhook received")

		// PRTG sends data using Content-Type: application/x-www-form-urlencoded
		err := r.ParseForm()
		if err != nil {
			msg := fmt.Sprintf("Error parsing form data: %s", err)
			log.Error().Err(err).Msg(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// webhook POST request from PRTG
		var data DataPRTG

		err = decoder.Decode(&data, r.PostForm)
		if err != nil {
			msg := fmt.Sprintf("Error parsing form data: %+v", err)
			log.Error().Err(err).Msg(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		kql, err := PRTGtoKQL(&data)
		if err != nil {
			msg := fmt.Sprintf("Failed to convert structure from PRTG to KQL: %+v", err)
			log.Error().Err(err).Msg(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		log.Debug().Msg("Processed PRTG form data")
		err = PopulateLAW(settings, kql)
		if err != nil {
			msg := fmt.Sprintf("Failed to populate the Log Analytics Workspace: %+v", err)
			log.Error().Err(err).Msg(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		log.Info().Str("Workspace ID", settings.WorkspaceID).
			Msg("Successfully forwarded the event to Log Analytics Workspace")

		fmt.Fprintf(w, "Ok")
	}
}

func Serve(settings *Settings) {
	decoder.IgnoreUnknownKeys(true)
	endpoint := strings.Trim(settings.Endpoint, "/")

	http.HandleFunc(fmt.Sprintf("/%s", endpoint), makePRTGHandler(settings))

	log.Info().Msgf("Forwarding PRTG Webhooks from %s:%d/%s to table '%s'", settings.Host, settings.Port, endpoint, settings.LogTable)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", settings.Host, settings.Port), nil); err != nil {
		log.Fatal().Msgf("Failed to start server: %v\n", err)
	}
}
