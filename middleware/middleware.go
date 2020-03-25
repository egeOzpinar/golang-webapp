package middleware

import(
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog"
	"net/http"
	"../sessions"
	"os"
)

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")
		_, ok := session.Values["user_id"] 
		if !ok {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
			http.Redirect(w, r, "/login", 302)
		return
		}
		handler.ServeHTTP(w, r)
	}
}