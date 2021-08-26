package main

import (
	"GPBChat/auth"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("views", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")

	if len(split) < 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	action := split[2]
	provider := split[3]

	switch action {
	case "login":
		//Login Action Handler
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s:%s", provider, err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		//Callback Action Handler
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to complete auth for %s : %s", provider, err), http.StatusBadRequest)
			return
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get user from %s : %s", provider, err), http.StatusInternalServerError)
			return
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get user from %s : %s", provider, err), http.StatusInternalServerError)
			return
		}
		authCookieValue := objx.New(map[string]interface{}{
			"name":       user.Name(),
			"email":      user.Email(),
			"avatar_url": user.AvatarURL(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the app.")
	flag.Parse()

	// setup gomniauth
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		facebook.New("1016005512494614", "afbf9e34ec777a3776252be92b84867b",
			"http://localhost:8080/auth/callback/facebook"),
		github.New("4fec0e0325d2a10685db", "faa5de782c74ee10d1a6cdd0f5eb5a0a55cc8e4b",
			"http://localhost:8080/auth/callback/github"),
		google.New("93978747738-eaar86h09qhl9cun8ha5314mm5t1ragi.apps.googleusercontent.com", "iEeXjIMdSQvVR-kWITWBxHy0",
			"http://localhost:8080/auth/callback/google"),
	)

	r := newRoom(UseGravatar)
	http.Handle("/chat", auth.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/room", r)

	go r.run()

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
