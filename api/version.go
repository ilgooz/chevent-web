package main

import (
	"net/http"

	"github.com/codeui/chevent-web/api/conf"
	"github.com/ilgooz/httpres"
)

var buildstamp string

type VersionResponse struct {
	Hash string `json:"hash"`
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	rp := VersionResponse{
		Hash: conf.Hash,
	}

	httpres.Json(w, http.StatusOK, rp)
}
