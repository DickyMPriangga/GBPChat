package main

import (
	"GPBChat/upload"
	"flag"
	"log"
	"net/http"
	"path/filepath"
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

var avatars Avatar = TryAvatar{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
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

	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.Handle("/upload", MustAuth(&templateHandler{filename: "upload.html"}))
	http.HandleFunc("/uploader", upload.UploadHandler)
	http.Handle("/avatars_img/", http.StripPrefix("/avatars_img/", http.FileServer(http.Dir("./avatars_img"))))
	http.Handle("/room", r)

	go r.run()

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
