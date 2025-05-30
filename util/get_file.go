package util

import (
    "os"
    "net/http"
    "crypto/tls"
    "log/slog"
    "io"
    "github.com/losenmann/iptv-toolkit/setup"
    "fmt"
    "regexp"
    "errors"
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

    data, err := io.ReadAll(f)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        return nil, err
    }
    return data, nil
}

func getFileRemote(url string) ([]byte, error) {
    tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
    client := &http.Client{Transport: tr}
    req, _ := http.NewRequest("GET", url, nil)
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        res, err = client.Get(url)
        if err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
            return nil, err
        }
    }
    defer res.Body.Close()

    data, err := io.ReadAll(res.Body)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        return nil, err
    }
    return data, nil
}