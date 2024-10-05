package convert

import (
    "iptv-toolkit/main/setup"
    "archive/zip"
    "encoding/xml"
    "encoding/binary"
    "compress/gzip"
    "fmt"
    "os"
    "log"
    "bytes"
    "regexp"
    "io/ioutil"
    "golang.org/x/text/encoding/charmap"
    "math"
	"time"
    "github.com/gabriel-vasile/mimetype"
    "log/slog"
)

var (
    v_epg_jtv = *setup.EpgPathDst + "/tvguide.zip"
    v_epg_xml = *setup.EpgPathDst + "/tvguide.xml"
    v_epg_xml_gz = *setup.EpgPathDst + "/tvguide.xml.gz"
)

type Tv struct {
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
    if err := os.MkdirAll(*setup.EpgPathDst, 0777); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    }
    switch mimetype.Detect(file).Extension() {
    case ".zip":
        err := ioutil.WriteFile(v_epg_jtv, file, 0644)
        if err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write " + v_epg_jtv)
            }
        }
        XmlToXmlGz(JtvToXml(v_epg_jtv))
    case ".xml":
        err := ioutil.WriteFile(v_epg_xml, file, 0644)
        if err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write " + v_epg_xml)
            }
        }
        XmlToXmlGz(file)
    case ".gz":
        err := ioutil.WriteFile(v_epg_xml_gz, file, 0644)
        if err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write " + v_epg_xml_gz)
            }
        }
        XmlGzToXml(file)
    default:
        if *setup.LogLVL <= 2 {
            slog.Warn("Unknown type epg")
        }
    }
}

func JtvToXml(file string) ([]byte) {
    var (
        epg Tv
        channel []Channel
        programme []Programme
        title []string
        duration []uint64
        pdt, ndx []byte
        channel_id int
        main_data int
    )

    r, err := zip.OpenReader(file)
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    for _, fpdt := range r.File {
        if re := regexp.MustCompile("pdt$").FindStringSubmatch(fpdt.Name); len(re) > 0 {
            channel_id++
            file_name := regexp.MustCompile(`\.pdt$`).ReplaceAllString(fpdt.Name, "")
            rc, err := fpdt.Open()
            if err != nil {
                log.Fatal(err)
            }

            pdt, _ = ioutil.ReadAll(rc)
            title = jtvParseTitle(pdt)
            rc.Close()
            for _, fndx := range r.File {
                if re2 := regexp.MustCompile("^" + file_name + ".ndx$").FindStringSubmatch(fndx.Name); len(re2) > 0 {
                    rc2, err := fndx.Open()
                    if err != nil {
                        log.Fatal(err)
                    }
                    ndx, _ = ioutil.ReadAll(rc2)
                    duration = jtvParseDuration(ndx)
                    rc2.Close()
                    break
                }
            }

            if len(title) <= len(duration) {
                main_data = len(title)
            } else {
                main_data = len(duration)
            }

            for k := 0; k < main_data; k++ {
                if k < main_data-1 {
                    programme = []Programme{Programme{Title:title[k],Channel:channel_id,Start:getTime(duration[k]),Stop:getTime(duration[k+1])}}
                } else {
                    programme = []Programme{Programme{Title:title[k],Channel:channel_id,Start:getTime(duration[k]),Stop:""}}
                }
                epg.Programme = append(epg.Programme, programme...)
                
            }
            channel = []Channel{Channel{DisplayName:file_name,Id:channel_id}}
            epg.Channel = append(epg.Channel, channel...)
        }
    }
    if data, err := xml.MarshalIndent(epg, "", "    "); err == nil {
        data = []byte(xml.Header + string(data))
        err := ioutil.WriteFile(v_epg_xml, data, 0644)
        if err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
        } else {
            if *setup.LogLVL <= 1 {
                slog.Info("Successfully write " + v_epg_xml)
            }
            return data
        }
    }
    return []byte{}
}

func XmlToXmlGz(data []byte) ([]byte) {
    var b bytes.Buffer
    gz := gzip.NewWriter(&b)
    gz.Write(data)
    gz.Close()
    err := ioutil.WriteFile(v_epg_xml_gz, b.Bytes(), 0644)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        return []byte{}
    } else {
        if *setup.LogLVL <= 1 {
            slog.Info("Successfully write " + v_epg_xml_gz)
        }
        return b.Bytes()
    }
    return []byte{}
}

func XmlGzToXml(data []byte) ([]byte) {
    reader := bytes.NewReader(data)
    gzreader, err := gzip.NewReader(reader);
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        return []byte{}
    }
    data, err = ioutil.ReadAll(gzreader);
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        return []byte{}
    }

    err = ioutil.WriteFile(v_epg_xml, data, 0644)
    if err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    } else {
        if *setup.LogLVL <= 1 {
            slog.Info("Successfully write " + v_epg_xml)
        }
        return data
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
        title_start, title_end, title_offset int
        offset_a [8]byte
        offset_b []byte
        dec = charmap.Windows1251.NewDecoder()
        jtv_headers1 = []byte("JTV 3.x TV Program Data\x0a\x0a\x0a")
        jtv_headers2 = []byte("JTV 3.x TV Program Data\xa0\xa0\xa0")
    )

    if bytes.Equal(pdt[0:26], jtv_headers1) || bytes.Equal(pdt[0:26], jtv_headers2) {
        data := pdt[26:]
        for title_start < len(data)-1 || title_end < len(data)-1 {
            offset_b, _ = dec.Bytes(data[title_start:title_start+1])
            copy(offset_a[8-len(offset_b):], offset_b)
            title_offset = int(binary.BigEndian.Uint64(offset_a[:])) + 2
            title_end = title_offset + title_start
            if title_end > len(data) {
                break
            } else {
                out, _ := dec.Bytes(data[title_start+2:title_end])
                title = append(title, string(out))
                title_start = title_start + title_offset
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
