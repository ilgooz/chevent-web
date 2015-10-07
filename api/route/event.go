package route

import (
	"log"
	"net/http"

	"github.com/codeui/chevent-web/api/ctx"
	"github.com/codeui/chevent-web/api/model"
	"github.com/ilgooz/formutils"
	"github.com/ilgooz/httpres"
	"github.com/ilgooz/paging"
)

type eventResponse struct {
	Event model.Event `json:"event"`
}

type eventsResponse struct {
	CurrentPage     int           `json:"current_page"`
	TotalPagesCount int           `json:"total_pages_counts"`
	Events          []model.Event `json:"events"`
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	event := model.Event{}

	if formutils.ParseSend(w, r, &event) {
		return
	}

	if err := ctx.M(r).DB("").C("events").Insert(&event); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	httpres.Json(w, http.StatusCreated, eventResponse{Event: event})
}

func ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	fields := listEventsForm{}

	if formutils.ParseSend(w, r, &fields) {
		return
	}

	q := ctx.M(r).DB("").C("events")

	totalCount, err := q.Count()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p := paging.Paging{
		Page:  fields.Page,
		Limit: fields.Limit,
		Count: totalCount,
	}.Calc()

	events := []model.Event{}

	if err := q.Find(nil).
		Skip(p.Offset).
		Limit(p.Limit).
		Sort("-date").
		All(&events); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rp := eventsResponse{
		CurrentPage:     p.Page,
		TotalPagesCount: p.TotalPages,
		Events:          events,
	}

	httpres.Json(w, http.StatusOK, rp)
}

type listEventsForm struct {
	Page  int `schema:"page"`
	Limit int `schema:"limit"`
}
