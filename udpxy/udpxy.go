package udpxy

import (
    "bufio"
    "github.com/gin-gonic/gin"
    "net"
    "strings"
)

const (
    maxDatagramSize = 32768
)

func Udpxy(c *gin.Context) {
    parts := strings.Split(c.Request.URL.Path, "/")
    if len(parts) < 3 || parts[1] != "udp" {
        c.String(400, "No address specified")
        return
    }

    addr, err := net.ResolveUDPAddr("udp", parts[2])
    if err != nil {
        c.String(500, err.Error())
        return
    }

    conn, err := net.ListenMulticastUDP("udp", nil, addr)
    if err != nil {
        c.String(500, err.Error())
        return
    }
    defer conn.Close()

    reader := bufio.NewReader(conn)
    b := make([]byte, maxDatagramSize)
    for {
        if n, err := reader.Read(b); err != nil {
            c.String(500, err.Error())
            return
        } else {
            c.Writer.Header().Set("Content-Type", "application/octet-stream")
            c.Writer.WriteHeader(200)

            if _, err := c.Writer.Write(b[:n]); err != nil {
                break
            }
        }
    }
}
