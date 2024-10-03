package convert

import (
    "bufio"
    "fmt"
    "encoding/xml"
    "os"
    "log/slog"
    "io/ioutil"
    "iptv-toolkit/main/setup"
    "regexp"
    "github.com/gabriel-vasile/mimetype"
    "bytes"
)

var (
    v_playlist_m3u = *setup.PlaylistPathDst + "/playlist.m3u"
    v_playlist_udpxy_m3u = *setup.PlaylistPathDst + "/playlist_udpxy.m3u"
    v_playlist_xml = *setup.PlaylistPathDst + "/playlist.xml"
    v_playlist_udpxy_xml = *setup.PlaylistPathDst + "/playlist_udpxy.xml"
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

func ConvertPlaylist(file []byte, udpxy string, epg string) {
    switch mimetype.Detect(file).Extension() {
    case ".m3u", ".m3u8":
        M3uToXml(file, udpxy)
    case ".xml":
        XmlToM3u(file, udpxy, epg)
    default:
        if *setup.LogLVL <= 2 {
            slog.Warn("Unknown type playlist")
        }
    }
}

func formatUdpxy(udpxy string ) string {
    return regexp.MustCompile("/[^/]*$").ReplaceAllString(udpxy, "") + "/udp/"
}

func XmlToM3u(file []byte, udpxy string, epg string,) {
    var data, data_udpxy string
    var playlist Playlist
    var track string = "\n#EXTINF:-1"
    var location string = ""
    var location_udpxy string = ""
    if udpxy != "" {
        udpxy = formatUdpxy(udpxy)
    }
    if epg != "" {
        data = "#EXTM3U url-tvg=\"" + epg + "\" m3uautoload=1 cache=500 deinterlace=1"
    } else {
        data = "#EXTM3U cache=500 deinterlace=1"
    }
    data_udpxy = data

    xml.Unmarshal(file, &playlist)
    for i := 0; i < len(playlist.TrackList.Track); i++ {
        track = "\n#EXTINF:-1"
        location = ""
        location_udpxy = ""

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
            if matched, _ := regexp.MatchString("^udp://@", playlist.TrackList.Track[i].Location); matched == true && udpxy != "" {
                location_udpxy = track + "\n" + regexp.MustCompile("^udp://@").ReplaceAllString(playlist.TrackList.Track[i].Location, udpxy)
            }
        } else {
            track = track + "\n"
        }
        data = data + location
        data_udpxy = data_udpxy + location_udpxy
    }

    // Write files XML
    err := ioutil.WriteFile(v_playlist_xml, file, 0644)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    } else {
        if *setup.LogLVL <= 1 {
            slog.Info("Successfully write playlist.xml")
        }
    }

    // Write files M3U
    f1, err := os.OpenFile(v_playlist_m3u, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    } else {
        f1.WriteString(data)
        defer f1.Close()
        if *setup.LogLVL <= 1 {
            slog.Info("Successfully write playlist.m3u")
        }
    }

    // Write files M3U Udpxy
    if matched, _ := regexp.MatchString("#EXTINF", data_udpxy); matched == true && udpxy != ""  {
        f2, err := os.OpenFile(v_playlist_udpxy_m3u, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
        if err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
            os.Exit(1)
        } else {
            f2.WriteString(data_udpxy)
            defer f2.Close()
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write playlist_udpxy.m3u")
            }
        }
    }
}

func M3uToXml(file []byte, udpxy string) {
    var playlist, playlist_udpxy Playlist
    var track []Track
    var image, title, psfile, location string
    var channel_id int
    if udpxy != "" {
        udpxy = formatUdpxy(udpxy)
    }

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
            if re := regexp.MustCompile("^udp://@").FindStringSubmatch(location); len(re) > 0 && udpxy != "" {
                track = []Track{Track{Channel_id:channel_id,Psfile:psfile,Image:image,Title:title,Location:regexp.MustCompile("^udp://@").ReplaceAllString(location, udpxy)}}
                playlist_udpxy.TrackList.Track = append(playlist_udpxy.TrackList.Track, track...)
            }
        }
    }

    if err := scanner.Err(); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    } else {
        // Write files M3U
        err := ioutil.WriteFile(v_playlist_m3u, file, 0644)
        if err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write playlist.m3u")
            }
        }

        // Write files XML
        if data, err := xml.MarshalIndent(playlist, "", "    "); err == nil {
		    data = []byte(xml.Header + string(data))
            err := ioutil.WriteFile(v_playlist_xml, data, 0644)
            if err != nil {
                if *setup.LogLVL <= 2 {
                    slog.Warn(fmt.Sprintf("%v", err))
                }
            } else {
                if *setup.LogLVL <= 1 {
                    slog.Info("Successfully write playlist.xml")
                }
            }
	    } else {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        }

        // Write files XML UDPXY
        if len(playlist_udpxy.TrackList.Track) > 0 && udpxy != ""  {
            if data, err := xml.MarshalIndent(playlist_udpxy, "", "    "); err == nil {
		        data = []byte(xml.Header + string(data))
                err := ioutil.WriteFile(v_playlist_udpxy_xml, data, 0644)
                if err != nil {
                    if *setup.LogLVL <= 2 {
                        slog.Warn(fmt.Sprintf("%v", err))
                    }
                } else {
                    if *setup.LogLVL <= 1 {
                        slog.Info("Successfully write playlist_udpxy.xml")
                    }
                }
	        } else {
                if *setup.LogLVL <= 2 {
                    slog.Warn(fmt.Sprintf("%v", err))
                }
            }
        }
    }
}