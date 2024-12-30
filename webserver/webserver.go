package webserver

import (
    "log/slog"
    "net/http"
    "strings"
    "strconv"
    "fmt"
    "iptv-toolkit/main/setup"
)

func Main(port int, path string) {
    mux := http.NewServeMux()
    mux.HandleFunc("/", pageRoot)
    mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(path))))
    
    if *setup.LogLVL <= 1 {
        slog.Info("starting web server on port: " + strconv.Itoa(port))
    }

    if err := http.ListenAndServe(":" + strconv.Itoa(port), mux); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    }
}


func pageRoot(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	} else {
        http.Redirect(w, r, "/files/", 301)
    }    
}


func hardened(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        index := strings.Index(r.URL.Path, "\x00")
        if index >= 0 {
            http.Error(w, "403 Forbidden", 403)
        } else {
            handler.ServeHTTP(w, r)
        }
    })
}
