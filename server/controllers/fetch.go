package controllers

import (
	"net/http"

	"github.com/secmohammed/reminders-cli-app/server/models"
	"github.com/secmohammed/reminders-cli-app/server/transport"
)

type fetcher interface {
	Fetch(ids []int) ([]models.Reminder, error)
}

func fetchReminders(service fetcher) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids, err := parseIDsParam(r.Context())
		if err != nil {
			transport.SendError(w, err)
			return
		}
		reminders, err := service.Fetch(ids)
		if err != nil {
			transport.SendError(w, err)
			return
		}
		transport.SendJSON(w, reminders, http.StatusOK)
	})
}
