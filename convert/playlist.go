package convert

import (
    "bufio"
    "bytes"
    "encoding/xml"
    "fmt"
    "github.com/gabriel-vasile/mimetype"
    "github.com/losenmann/iptv-toolkit/setup"
    "log/slog"
    "os"
    "regexp"
)

type Playlist struct {
    XMLName xml.Name `xml:"playlist"`
    Title string `xml:"title"`
    TrackList TrackList `xml:"trackList"`
}

type TrackList struct {
    Track []Track `xml:"track"`
}

type Track struct {
    Channel_id int `xml:"channel_id"`
    Location string `xml:"location"`
    Title string `xml:"title"`
    Image string `xml:"image"`
    Psfile string `xml:"psfile"`
    Zoom string `xml:"zoom"`
    Is_external string `xml:"is_external"`
}

func ConvertPlaylist(file []byte, udpxy, epg bool) {
    var (
        path_playlist_m3u = *setup.PlaylistDir + "/playlist.m3u"
        path_playlist_udpxy_m3u = *setup.PlaylistDir + "/playlist_udpxy.m3u"
        path_playlist_xml = *setup.PlaylistDir + "/playlist.xml"
        path_playlist_udpxy_xml = *setup.PlaylistDir + "/playlist_udpxy.xml"
        playlist_xml, playlist_m3u, playlist_udpxy_xml, playlist_udpxy_m3u []byte
    )

    switch mimetype.Detect(file).Extension() {
    case ".m3u", ".m3u8":
        playlist_m3u = file
        playlist_xml = M3uToXml(file)
        if udpxy {
            if epg {
                playlist_udpxy_m3u = UdpToUdpxy(file, *setup.Address)
                playlist_udpxy_xml = UdpToUdpxy(playlist_xml, *setup.Address)
            } else {
                playlist_udpxy_m3u = UdpToUdpxy(file, "")
                playlist_udpxy_xml = UdpToUdpxy(playlist_xml, "")
            }
        }
    case ".xml":
        playlist_xml = file
        if epg {
            playlist_m3u = XmlToM3u(file, *setup.Address)
        } else {
            playlist_m3u = XmlToM3u(file, "")
        }
        if udpxy {
            if epg {
                playlist_udpxy_xml = UdpToUdpxy(file, *setup.Address)
                playlist_udpxy_m3u = UdpToUdpxy(playlist_m3u, *setup.Address)
            } else {
                playlist_udpxy_xml = UdpToUdpxy(file, "")
                playlist_udpxy_m3u = UdpToUdpxy(playlist_m3u, "")
            }
        }
    default:
        if *setup.LogLVL <= 2 {
            slog.Warn("Unknown type playlist")
        }
    }

    if err := os.MkdirAll(*setup.PlaylistDir, 0777); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    } else {
        if err := os.WriteFile(path_playlist_m3u, playlist_m3u, 0644); err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write " + path_playlist_m3u)
            }
        }
        if err := os.WriteFile(path_playlist_xml, playlist_xml, 0644); err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write " + path_playlist_xml)
            }
        }
        if len(playlist_udpxy_m3u) > 0 {
            if err := os.WriteFile(path_playlist_udpxy_m3u, playlist_udpxy_m3u, 0644); err != nil {
                if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
                }
            } else {
                if *setup.LogLVL <= 1 {
                    slog.Info("Successfully write " + path_playlist_udpxy_m3u)
                }
            }
        }
        if len(playlist_udpxy_xml) > 0 {
            if err := os.WriteFile(path_playlist_udpxy_xml, playlist_udpxy_xml, 0644); err != nil {
                if *setup.LogLVL <= 2 {
                    slog.Warn(fmt.Sprintf("%v", err))
                }
            } else {
                if *setup.LogLVL <= 1 {
                    slog.Info("Successfully write " + path_playlist_udpxy_xml)
                }
            }
        }
    }
}

func formatUdpxy(udpxy string) string {
    return regexp.MustCompile("/[^/]*$").ReplaceAllString(udpxy, "") + "/udp/"
}

func UdpToUdpxy(file []byte, udpxy string) ([]byte) {
    return []byte(regexp.MustCompile(`udp://@`).ReplaceAllString(string(file), formatUdpxy(udpxy)))
}

func XmlToM3u(file []byte, epg ...string) ([]byte) {
    var data string
    var playlist Playlist
    var track string = "\n#EXTINF:-1"
    var location string = ""

    if len(epg) > 0 && epg[0] != "" {
        data = "#EXTM3U url-tvg=\"" + epg[0] + "\" m3uautoload=1 cache=500 deinterlace=1"
    } else {
        data = "#EXTM3U cache=500 deinterlace=1"
    }

    xml.Unmarshal(file, &playlist)
    for i := 0; i < len(playlist.TrackList.Track); i++ {
        track = "\n#EXTINF:-1"
        location = ""

        if playlist.TrackList.Track[i].Psfile != "" {
            track = track + ",tvg-name=\"" + playlist.TrackList.Track[i].Psfile + "\""
        }
        if playlist.TrackList.Track[i].Image != "" {
            track = track + ",tvg-logo=\"" + playlist.TrackList.Track[i].Image + "\""
        }
        if playlist.TrackList.Track[i].Title != "" {
            track = track + "," + playlist.TrackList.Track[i].Title
        }
        if playlist.TrackList.Track[i].Location != "" {
            location = track + "\n" + playlist.TrackList.Track[i].Location
        } else {
            track += "\n"
        }
        data = data + location
    }

    return []byte(data)
}

func M3uToXml(file []byte) ([]byte) {
    var playlist Playlist
    var track []Track
    var image, title, psfile, location string
    var channel_id int

    scanner := bufio.NewScanner(bytes.NewReader(file))
    for scanner.Scan() {
        if extinf := regexp.MustCompile(`^#EXTINF`).FindStringSubmatch(scanner.Text()); len(extinf) > 0 {
            channel_id++
            image = ""
            title = ""
            psfile = ""
            if tvg_name := regexp.MustCompile(`tvg-name="[^"]*`).FindStringSubmatch(scanner.Text()); len(tvg_name) > 0 {
                psfile = tvg_name[0][10:]
            }
            if tvg_logo := regexp.MustCompile(`tvg-logo="[^"]*`).FindStringSubmatch(scanner.Text()); len(tvg_logo) > 0 {
                image = tvg_logo[0][10:]
            }
            if tvg_title := regexp.MustCompile(`[^,]*$`).FindStringSubmatch(scanner.Text()); len(tvg_title) > 0 {
                title = tvg_title[0]
            }
        }
        if tvg_location := regexp.MustCompile(`^[^#]*$`).FindStringSubmatch(scanner.Text()); len(tvg_location) > 0 {
            location = tvg_location[0]
            track = []Track{Track{Channel_id:channel_id,Psfile:psfile,Image:image,Title:title,Location:location}}
            playlist.TrackList.Track = append(playlist.TrackList.Track, track...)
        }
    }

    if err := scanner.Err(); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    } else {
        if data, err := xml.MarshalIndent(playlist, "", "    "); err == nil {
            return []byte(xml.Header + string(data))
        }
    }
    return []byte{}
}
