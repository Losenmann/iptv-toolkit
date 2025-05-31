package convert

import (
    "github.com/losenmann/iptv-toolkit/setup"
    "archive/zip"
    "encoding/xml"
    "encoding/binary"
    "compress/gzip"
    
    "fmt"
    "os"
    "bytes"
    "regexp"
    "io"
    "golang.org/x/text/encoding/charmap"
    "math"
    "time"
    "github.com/gabriel-vasile/mimetype"
    "log/slog"
    "strconv"
)

type Epg struct {
    XMLName xml.Name `xml:"tv"`
    Channel []Channel `xml:"channel"`
    Programme []Programme `xml:"programme"`
}
type Channel struct {
    DisplayName string `xml:"display-name"`
    Icon Icon `xml:"icon"`
    Id int `xml:"id,attr"`
}
type Icon struct {
    Src string `xml:"src,attr"`
}
type Programme struct {
    Title string `xml:"title"`
    Start string `xml:"start,attr"`
    Stop string `xml:"stop,attr"`
    Channel int `xml:"channel,attr"`
}

func ConvertEpg(file []byte) {
    var (
        path_epg_jtv = *setup.EpgDir + "/epg.zip"
        path_epg_xml = *setup.EpgDir + "/epg.xml"
        path_epg_xmlgz = *setup.EpgDir + "/epg.xml.gz"
        epg_jtv, epg_xml, epg_xmlgz []byte
    )
    
    switch mimetype.Detect(file).Extension() {
    case ".zip":
        epg_jtv = file
        epg_xml = JtvToXml(file)
        epg_xmlgz = XmlToXmlGz(epg_xml)
    case ".xml":
        epg_xml = file
        epg_xmlgz = XmlToXmlGz(file)
        epg_jtv = XmlToJtv(file)
    case ".gz":
        epg_xmlgz = file
        epg_xml = XmlGzToXml(file)
        epg_jtv = XmlToJtv(epg_xml)
    default:
        if *setup.LogLVL <= 2 {
            slog.Warn("Unknown type epg")
        }
    }
    if err := os.MkdirAll(*setup.EpgDir, 0777); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    } else {
        if err := os.WriteFile(path_epg_jtv, epg_jtv, 0644); err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write " + path_epg_jtv)
            }
        }
        if err := os.WriteFile(path_epg_xml, epg_xml, 0644); err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write " + path_epg_xml)
            }
        }
        if err := os.WriteFile(path_epg_xmlgz, epg_xmlgz, 0644); err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write " + path_epg_xmlgz)
            }
        }
    }
}

func JtvToXml(file []byte) ([]byte) {
    var (
        epg Epg
        channel []Channel
        programme []Programme
        title []string
        duration []uint64
        channel_id, main_data int
    )

    reader := bytes.NewReader(file)
    zipreader, err := zip.NewReader(reader, int64(len(file)))
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        return []byte{}
    }

    for _, fpdt := range zipreader.File {
        if ch_name := regexp.MustCompile(`\.pdt$`).ReplaceAllString(fpdt.Name, ""); ch_name != fpdt.Name {
            if rc, err := zipreader.Open(ch_name + ".pdt"); err != nil {
                if *setup.LogLVL <= 2 {
                    slog.Warn(fmt.Sprintf("%v", err))
                }
                continue
            } else {
                pdt, _ := io.ReadAll(rc)
                title = jtvParseTitle(pdt)
                rc.Close()
            }
            if rc, err := zipreader.Open(ch_name + ".ndx"); err != nil {
                if *setup.LogLVL <= 2 {
                    slog.Warn(fmt.Sprintf("%v", err))
                }
                continue
            } else {
                ndx, _ := io.ReadAll(rc)
                duration = jtvParseDuration(ndx)
                rc.Close()
            }

            if len(title) <= len(duration) {
                main_data = len(title)
            } else {
                main_data = len(duration)
            }
            
            if main_data > 0 {
                channel_id++
                for k := 0; k < main_data; k++ {
                    if k < main_data-1 {
                        programme = []Programme{Programme{Title:title[k],Channel:channel_id,Start:getTime(duration[k]),Stop:getTime(duration[k+1])}}
                    } else {
                        programme = []Programme{Programme{Title:title[k],Channel:channel_id,Start:getTime(duration[k]),Stop:""}}
                    }
                    epg.Programme = append(epg.Programme, programme...)
                }
                channel = []Channel{Channel{DisplayName:ch_name,Id:channel_id}}
                epg.Channel = append(epg.Channel, channel...)
            }
        }
    }

    if data, err := xml.MarshalIndent(epg, "", "    "); err == nil {
        return []byte(xml.Header + string(data))
    }
    return []byte{}
}

func XmlToJtv(data []byte) ([]byte) {
    var epg Epg
    xml.Unmarshal(data, &epg)
    return jtvCreateFileFromXml(epg)
}

func XmlToXmlGz(data []byte) ([]byte) {
    var buff bytes.Buffer
    gz := gzip.NewWriter(&buff)
    gz.Write(data)
    gz.Close()
    return buff.Bytes()
}

func XmlGzToXml(data []byte) ([]byte) {
    reader := bytes.NewReader(data)
    if gzreader, err := gzip.NewReader(reader); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    } else {
        if data, err := io.ReadAll(gzreader); err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            return data
        }
    }
    return []byte{}
}

func getTime(data uint64) string {
    maxd := time.Duration(math.MaxInt64).Truncate(100 * time.Nanosecond)
    maxdUnits := uint64(maxd / 100) // number of 100-ns units

    t := time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)
    for data > maxdUnits {
        t = t.Add(maxd)
        data -= maxdUnits
    }
    if data != 0 {
        t = t.Add(time.Duration(data * 100))
    }
    return t.Format("20060102150405")
}

func jtvParseTitle(pdt []byte) ([]string) {
    var (
        title []string
        title_start, title_end, title_offset, title_len int
        offset_a [8]byte
        offset_b []byte
        dec = charmap.Windows1251.NewDecoder()
        jtv_headers = [][]byte{[]byte("JTV 3.x TV Program Data\x0a\x0a\x0a"), []byte("JTV 3.x TV Program Data\xa0\xa0\xa0")}
    )

    if bytes.Equal(pdt[0:26], jtv_headers[0]) || bytes.Equal(pdt[0:26], jtv_headers[1]) {
        data := pdt[26:]
        for title_start < len(data)-1 && title_end < len(data)-1 {
            title_start = title_end
            title_len = title_start+2
            offset_b = data[title_start:title_len-1]
            copy(offset_a[8-len(offset_b):], offset_b)
            title_offset = int(binary.BigEndian.Uint64(offset_a[:])) + 2
            title_end = title_start + title_offset
            if title_end > len(data) {
                break
            } else {
                out, _ := dec.Bytes(data[title_len:title_end])
                if len(out) > 0 {
                    title = append(title, string(out))
                }
            }
        }
    }
    return title
}

func jtvParseDuration(ndx []byte) ([]uint64) {
    var start_duration int = 0
    var end_duration int = 12
    var duration []uint64
    data := ndx[2:]
    for start_duration < len(data) {
        duration = append(duration, uint64(binary.LittleEndian.Uint64(data[start_duration+2:end_duration-2])))
        start_duration = start_duration + 12
        end_duration = end_duration + 12
    }
    return duration
}

func jtvCreateFileFromXml(epg Epg) ([]byte) {
    var (
        buff bytes.Buffer
        jtv_headers = [][]byte{[]byte("JTV 3.x TV Program Data\x0a\x0a\x0a"), []byte("JTV 3.x TV Program Data\xa0\xa0\xa0")}
        enc = charmap.Windows1251.NewEncoder()
        filetime uint64 = 116444736000000000
    )

    zipW := zip.NewWriter(&buff)
    for _, ch := range(epg.Channel) {
        if id, _ := strconv.Atoi(ch.DisplayName); id <= 1154{
            pdt := jtv_headers[0]
            var ndx []byte
            var count int
            for _, prog := range(epg.Programme) {
                if prog.Channel == ch.Id {
                    count++
                    prog_len := make([]byte, 8)
                    prog_title, _ := enc.Bytes([]byte(prog.Title))
                    binary.LittleEndian.PutUint64(prog_len, uint64(len(prog_title)))
                    pdt = append(pdt, append(prog_len[:2], prog_title...)...)

                    duration := make([]byte, 8)
                    t, _ := time.Parse("20060102150405", prog.Start)
                    binary.LittleEndian.PutUint64(duration, uint64((t.Unix()*10000000))+filetime)
                    ndx = append(ndx, append([]byte("\x00\x00"), append(duration, prog_len[:2]...)...)...)
                }
            }
            ndx_cnt := make([]byte, 8)
            binary.LittleEndian.PutUint64(ndx_cnt, uint64(count))
            ndx = append(ndx_cnt[:2], ndx...) 

            if f, err := zipW.Create(ch.DisplayName + ".pdt"); err != nil {
                if *setup.LogLVL <= 2 {
                    slog.Warn(fmt.Sprintf("%v", err))
                }
                continue
            } else {
                if _, err := f.Write(pdt); err != nil {
                    if *setup.LogLVL <= 2 {
                        slog.Warn(fmt.Sprintf("%v", err))
                    }
                    continue
                }
                if f, err := zipW.Create(ch.DisplayName + ".ndx"); err != nil {
                    if *setup.LogLVL <= 2 {
                        slog.Warn(fmt.Sprintf("%v", err))
                    }
                    continue
                } else {
                    if _, err := f.Write(ndx); err != nil {
                        if *setup.LogLVL <= 2 {
                            slog.Warn(fmt.Sprintf("%v", err))
                        }
                        continue
                    }
                }
            }
        }
    }

    err := zipW.Close()
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    } else {
        return buff.Bytes()
    }
    return []byte{}
}
