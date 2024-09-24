package main

import (
//    "os"
//    "fmt"
//    "bufio"
//    "compress/gzip"
//    "bytes"
//    "io/ioutil"
    "log/slog"
    "iptv-toolkit/main/setup"
    "iptv-toolkit/main/convert"
    "iptv-toolkit/main/util"
//	"errors"
)

func init() {
    setup.Initgo()
}

func main() {
    playlist, err := util.GetFile(*setup.Playlist)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn("playlist file not specified")
        }
    } else {
        convert.ConvertPlaylist(playlist, *setup.PlalistUdpxy, *setup.EmbedEPG)
    }
/*
    epg, err := util.GetFile(*setup.Epg)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn("epg file not specified")
        }
    } else {
        convert.ConvertEpg(epg)
    }
*/
}
