package petservice

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

type contextKeyID int

const (
	keyID contextKeyID = iota
)

func keyExtractorMiddleware(urlParam string, id contextKeyID) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			up := chi.URLParam(r, urlParam)
			ctx := context.WithValue(r.Context(), id, up)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AttachPetserviceEndpoints attaches the endpoint for a pet service to the given router
func AttachPetserviceEndpoints(r chi.Router, s *PetService) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/", makeHandler(s.PostPet))
		r.Route("/{key}", func(r chi.Router) {
			r.Use(keyExtractorMiddleware("key", keyID))
			r.Get("/", makeHandler(s.GetPet))
			r.Put("/", makeHandler(s.PutPet))
			r.Delete("/", makeHandler(s.DeletePet))
		})
	})
}

func renderErrorResponse(w http.ResponseWriter, r *http.Request, e error) {
	log.Infof("Error serving request. %v", e.Error())
	if apiError, ok := e.(APIError); ok {
		http.Error(w, apiError.ResponseMessage, apiError.StatusCode)
	} else {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(serveRequest func(*http.Request) (int, interface{}, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode, response, err := serveRequest(r)
		if err != nil {
			renderErrorResponse(w, r, err)
			return
		}
		render.Status(r, statusCode)
		render.JSON(w, r, response)
	})
}
