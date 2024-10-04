package setup

import (
    "os"
    "flag"
    "regexp"
    "strconv"
)

var (
    LogLVL              = flag.Int("l", 1, "Log lvl")
    Epg                 = flag.String("e", "", "EPG URI")
    EpgPathDst          = flag.String("E", "./files/tvguide", "Path EPGs storage")
    Playlist            = flag.String("p", "", "Playlist URI")
    PlaylistPathDst     = flag.String("P", "./files/playlist", "Path playlists storage")
    WebPort             = flag.Int("w", 4023, "Web Server port")
    WebPath             = flag.String("f", "./files", "Web Server path")
    EmbedEPG            = flag.String("i", "", "Embed a link to EPG in the playlist")
    PlalistUdpxy        = flag.String("u", "", "Create a playlist with embedded udpxy")
    Health              = flag.Bool("H", false, "Start Healthcheck server")
    Udpxy               = flag.Bool("U", false, "Start UDPXY proxy")
    Schedule            = flag.Bool("S", false, "Start Schedule job")
    Web                 = flag.Bool("W", false, "Start Web Server")
    Crontab             = flag.String("c", "30 6 * * *", "Create a playlist with embedded udpxy")
)

func Initgo() {
    flag.Parse()
    if os.Getenv("IPTVTOOLKIT_WEB_PATH") != "" && *WebPath == "./files" {
        *WebPath = os.Getenv("IPTVTOOLKIT_WEB_PATH")
    }
    if os.Getenv("IPTVTOOLKIT_WEB_PORT") != "" && *WebPort == 4023 {
        *WebPort, _ = strconv.Atoi(os.Getenv("IPTVTOOLKIT_WEB_PORT"))
    }
    if os.Getenv("IPTVTOOLKIT_EPG") != "" && *Epg == "" {
        *Epg = os.Getenv("IPTVTOOLKIT_EPG")
    }
    if os.Getenv("IPTVTOOLKIT_PLAYLIST") != "" && *Playlist == "" {
        *Playlist = os.Getenv("IPTVTOOLKIT_PLAYLIST")
    }
    if os.Getenv("IPTVTOOLKIT_EPG_DST") != "" && *EpgPathDst == "./files/tvguide" {
        *EpgPathDst = os.Getenv("IPTVTOOLKIT_EPG_DST")
    }
    if os.Getenv("IPTVTOOLKIT_PLAYLIST_DST") != "" && *PlaylistPathDst == "./files/playlist" {
        *PlaylistPathDst = os.Getenv("IPTVTOOLKIT_PLAYLIST_DST")
    }
    if os.Getenv("IPTVTOOLKIT_PLAYLIST_UDPXY") != "" && *PlalistUdpxy == "" {
        *PlalistUdpxy = os.Getenv("IPTVTOOLKIT_PLAYLIST_UDPXY")
    }
    if os.Getenv("IPTVTOOLKIT_PLAYLIST_EMBED_EPG") != "" && *EmbedEPG == "" {
        *EmbedEPG = os.Getenv("IPTVTOOLKIT_PLAYLIST_EMBED_EPG")
    }
    if os.Getenv("IPTVTOOLKIT_CRONTAB") != "" && *Crontab == "30 6 * * *" {
        *Crontab = os.Getenv("IPTVTOOLKIT_CRONTAB")
    }

    *EpgPathDst = regexp.MustCompile("/$").ReplaceAllString(*EpgPathDst, "")
    *PlaylistPathDst = regexp.MustCompile("/$").ReplaceAllString(*PlaylistPathDst, "")
}
