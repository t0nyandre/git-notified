package github

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (p *Provider) Login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, p.AuthURL, http.StatusMovedPermanently)
}

func (p *Provider) Callback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User: %v", "t0nyandre")
}

func (p *Provider) Routes(r *chi.Mux) chi.Router {
	r.Get("/github/login", p.Login)
	r.Get("/github/callback", p.Callback)
	return r
}
