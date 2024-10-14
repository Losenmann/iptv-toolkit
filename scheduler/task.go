package scheduler

import (
    "log/slog"
    "iptv-toolkit/main/setup"
    "iptv-toolkit/main/util"
    "iptv-toolkit/main/convert"
)

func Task() {
    //convert.XmlToJtv()
    if *setup.Playlist != "" {
        if playlist, err := util.GetFile(*setup.Playlist); err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn("playlist file not specified")
            }
        } else {
            convert.ConvertPlaylist(playlist, *setup.PlalistUdpxy, *setup.EmbedEPG)
        }
    }

    if *setup.Epg != "" {
        if epg, err := util.GetFile(*setup.Epg); err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn("epg file not specified")
            }
        } else {
            convert.ConvertEpg(epg)
        }
    }
}