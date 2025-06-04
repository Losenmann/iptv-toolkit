package setup

import (
    "flag"
    "os"
    "regexp"
    "strconv"
)

var (
    LogLVL              = flag.Int("l", 1, "Log lvl")
    Address             = flag.String("a", getEnv("IPTVTOOLKIT_ADDRESS", ""), "EPG URI")
    Epg                 = flag.String("e", getEnv("IPTVTOOLKIT_EPG", ""), "EPG URI")
    EpgDir              = flag.String("E", strippath(getEnv("IPTVTOOLKIT_EPG_DIR", "./tvguide")), "Path EPGs storage")
    Playlist            = flag.String("p", getEnv("IPTVTOOLKIT_PLAYLIST", ""), "Playlist URI")
    PlaylistDir         = flag.String("P", strippath(getEnv("IPTVTOOLKIT_PLAYLIST_DIR", "./playlist")), "Path playlists storage")
    WebPathUdpxy        = flag.String("u", getEnv("IPTVTOOLKIT_WEB_PATH_UDXPY", "/udp/:addr"), "Web Server multicast path")
    WebPathFiles        = flag.String("f", getEnv("IPTVTOOLKIT_WEB_PATH_FILES", "/files"), "Web Server static files path")
    WebDir              = flag.String("W", getEnv("IPTVTOOLKIT_WEB_DIR", "./files"), "Directory on the host for web server files")
    WebPort             = flag.Int("w", aToi(getEnv("IPTVTOOLKIT_WEB_PORT", "4022")), "Web Server port")
    Crontab             = flag.String("c", getEnv("IPTVTOOLKIT_CRONTAB", "30 6 * * *"), "Сrontab style task schedule")
    EmbedEPG            = flag.Bool("t", aTob(getEnv("IPTVTOOLKIT_EMBED_EPG", "false"), false), "UDPXY URI")
    EmbedUdpxy          = flag.Bool("d", aTob(getEnv("IPTVTOOLKIT_EMBED_UDPXY", "false"), false), "Create a playlist with embedded udpxy")
    Udpxy               = flag.Bool("U", false, "Start UDPXY proxy")
    Files               = flag.Bool("F", false, "Start File Server")
    Schedule            = flag.Bool("S", false, "Start Schedule job")
    Health              = flag.Bool("H", false, "Start Healthcheck endpoint")
    Metric              = flag.Bool("M", false, "Start Metrics endpoint")
)

func getEnv(key, fallback string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        return fallback
    }
    return value
}

func aToi(value string) int {
    a, _ := strconv.Atoi(value)
    return a
}

func aTob(s string, fallback bool) bool {
    if value, err := strconv.ParseBool(s); err != nil {
        return fallback
    } else {
        return value
    }
}

func Initgo() {
    flag.Parse()
    *EpgDir = regexp.MustCompile(`/$`).ReplaceAllString(*EpgDir, "")
    *PlaylistDir = regexp.MustCompile(`/$`).ReplaceAllString(*PlaylistDir, "")
}

func strippath(a string) string {
    return regexp.MustCompile(`/$`).ReplaceAllString(a, "")
}
