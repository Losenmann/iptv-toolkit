package setup

import (
    "flag"
    "os"
    "regexp"
)

var (
    LogLVL              = flag.Int("l", 1, "Log lvl")
    Address             = flag.String("a", getEnv("IPTVTOOLKIT_ADDRESS", "").(string), "EPG URI")
    Epg                 = flag.String("e", getEnv("IPTVTOOLKIT_EPG", "").(string), "EPG URI")
    EpgDir              = flag.String("E", strippath(getEnv("IPTVTOOLKIT_EPG_DIR", "./files/tvguide").(string)), "Path EPGs storage")
    Playlist            = flag.String("p", getEnv("IPTVTOOLKIT_PLAYLIST", "").(string), "Playlist URI")
    PlaylistDir         = flag.String("P", strippath(getEnv("IPTVTOOLKIT_PLAYLIST_DIR", "./files/playlist").(string)), "Path playlists storage")
    WebPathUdpxy        = flag.String("u", getEnv("IPTVTOOLKIT_WEB_PATH_UDXPY", "/udp/:addr").(string), "Web Server multicast path")
    WebPathFiles        = flag.String("f", getEnv("IPTVTOOLKIT_WEB_PATH_FILES", "/files").(string), "Web Server static files path")
    WebDir              = flag.String("W", getEnv("IPTVTOOLKIT_WEB_DIR", "./files").(string), "Directory on the host for web server files")
    WebPort             = flag.Int("w", getEnv("IPTVTOOLKIT_WEB_PORT", "4022").(int), "Web Server port")
    Crontab             = flag.String("c", getEnv("IPTVTOOLKIT_CRONTAB", "30 6 * * *").(string), "Сrontab style task schedule")
    EmbedEPG            = flag.Bool("t", getEnv("IPTVTOOLKIT_EMBED_EPG", "false").(bool), "UDPXY URI")
    EmbedUdpxy          = flag.Bool("d", getEnv("IPTVTOOLKIT_EMBED_UDPXY", "false").(bool), "Create a playlist with embedded udpxy")
    Udpxy               = flag.Bool("U", false, "Start UDPXY proxy")
    Files               = flag.Bool("F", false, "Start File Server")
    Schedule            = flag.Bool("S", false, "Start Schedule job")
    Health              = flag.Bool("H", false, "Start Healthcheck endpoint")
    Metric              = flag.Bool("M", false, "Start Metrics endpoint")
)

func getEnv(key, fallback string) any {
    value, exists := os.LookupEnv(key)
    if !exists {
        return fallback
    }
    return value
}

func Initgo() {
    flag.Parse()
    *EpgDir = regexp.MustCompile(`/$`).ReplaceAllString(*EpgDir, "")
    *PlaylistDir = regexp.MustCompile(`/$`).ReplaceAllString(*PlaylistDir, "")
}

func strippath(a string) string {
    return regexp.MustCompile(`/$`).ReplaceAllString(a, "")
}
