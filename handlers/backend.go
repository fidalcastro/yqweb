package handlers

import (
	"net/http"

	"github.com/fidalcastro/yqweb/views"
)

func GetFoo(w http.ResponseWriter, r *http.Request) {
	views.Index().Render(r.Context(), w)
}

func PostFoo(w http.ResponseWriter, r *http.Request) {
	views.FilterResult("Answer").Render(r.Context(), w)
}

func GetFilter(w http.ResponseWriter, r *http.Request) {
	views.Filter().Render(r.Context(), w)
}
