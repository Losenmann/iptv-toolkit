package util

import (
    "os"
    "net/http"
    "log/slog"
    "io/ioutil"
    "iptv-toolkit/main/setup"
    "fmt"
    "regexp"
    "errors"
)

var (
    v_tvguide = *setup.EpgPathDst + "/tvguide"
    v_playlist = *setup.PlaylistPathDst + "/playlist"
)

func GetFile(uri string) ([]byte, error) {
    if uri != "" {
        re := regexp.MustCompile("^(http[s]*://|./|/|[^(http://|https://)])").FindStringSubmatch(uri)
        switch regexp.MustCompile("^(./|/|[^(http://|https://)])").ReplaceAllString(re[0], "local") {
            case "http://", "https://":
                return getFileRemote(uri)
            case "local":
                return getFileLocal(uri)
        }
    }
    return nil, errors.New("File not specified")
}

func getFileLocal(file string) ([]byte, error) {
    f, err := os.Open(file)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        return nil, err
    }
    defer f.Close()

    data, err := ioutil.ReadAll(f)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        return nil, err
    }
    return data, nil
}

func getFileRemote(url string) ([]byte, error) {
    req, _ := http.NewRequest("GET", url, nil)
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        return nil, err
    }
    defer res.Body.Close()

    data, err := ioutil.ReadAll(res.Body)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        return nil, err
    }
    return data, nil
}
