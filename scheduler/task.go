package scheduler

import (
    "github.com/losenmann/iptv-toolkit/convert"
    "github.com/losenmann/iptv-toolkit/setup"
    "github.com/losenmann/iptv-toolkit/util"
    "log/slog"
)

func Task() {
    if *setup.Playlist != "" {
        if playlist, err := util.GetFile(*setup.Playlist); err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn("playlist file not specified")
            }
        } else {
            convert.ConvertPlaylist(playlist, *setup.PlaylistEmbedUdpxy, *setup.PlaylistEmbedEPG)
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
