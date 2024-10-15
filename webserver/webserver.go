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
    fs := http.FileServer(http.Dir(path))
    mux.Handle("/files/", http.StripPrefix("/files", fs))
    mux.Handle("/static/", http.StripPrefix("/static", fs))

    if *setup.LogLVL <= 1 {
        slog.Info("starting web server on port: " + strconv.Itoa(port))
    }
    if err := http.ListenAndServe(":" + strconv.Itoa(port), mux); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
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