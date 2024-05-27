package github_handler

import (
	"context"
	"net/http"

	"github.com/DervexDev/ghloc/src/server/rest"
	"github.com/DervexDev/ghloc/src/service/github_stat"
	"github.com/DervexDev/ghloc/src/service/loc_count"
	"github.com/go-chi/chi/v5"
)

type Service interface {
	GetStat(ctx context.Context, user, repo, branch, token string, filter, matcher *string, noLOCProvider bool, tempStorage github_stat.TempStorage) (*loc_count.StatTree, error)
}

type GetStatHandler struct {
	Service    Service
}

func (h *GetStatHandler) RegisterOn(router chi.Router) {
	router.Get("/{user}/{repo}/{branch}", h.ServeHTTP)
}

func (h GetStatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	repo := chi.URLParam(r, "repo")
	branch := chi.URLParam(r, "branch")
	token := r.Header.Get("Authorization")

	r.ParseForm()

	noLOCProvider := false
	tempStorage := github_stat.TempStorageFile

	var filter *string
	if filters := r.Form["filter"]; len(filters) >= 1 {
		filter = &filters[0]
	}

	var matcher *string
	if matchers := r.Form["match"]; len(matchers) >= 1 {
		matcher = &matchers[0]
	}

	stat, err := h.Service.GetStat(r.Context(), user, repo, branch, token, filter, matcher, noLOCProvider, tempStorage)
	if err != nil {
		rest.WriteResponse(w, r, err, true)
		return
	}
	w.Header().Add("Cache-Control", "public, max-age=300")
	rest.WriteResponse(w, r, (*rest.SortedStat)(stat), r.FormValue("pretty") != "false")
}