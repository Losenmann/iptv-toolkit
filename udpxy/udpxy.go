package udpxy

import (
	"iptv-toolkit/main/setup"
	"encoding/hex"
	"log/slog"
    "fmt"
	"log"
	"net"
	"time"
	"os/exec"
	"os"
)

var (
    v_tvguide = *setup.EpgPathDst + "/tvguide"
    v_playlist = *setup.PlaylistPathDst + "/playlist"
)

const (
//	srvAddr         = "224.0.0.1:9999"
	srvAddr         = "239.255.43.72:1234"
	maxDatagramSize = 8192
)

func Main() {
	go ping(srvAddr)
	serveMulticastUDP(srvAddr, msgHandler)
}

func ping(a string) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	_, err = net.DialUDP("udp", nil, addr)
	for {
	//	c.Write([]byte("hello, world\n"))
		time.Sleep(1 * time.Second)
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "bytes read from", src)
	log.Println(hex.Dump(b[:n]))
	log.Println(b[:n])
}

func serveMulticastUDP(a string, h func(*net.UDPAddr, int, []byte)) {
	

	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.ListenMulticastUDP("udp", nil, addr)
	l.SetReadBuffer(maxDatagramSize)
	
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		h(src, n, b)
	}
}

func UdpxyExt() {
	cmd := exec.Command("udpxy", "-p", "4022", "-vTSl", "/proc/1/fd/1")
	_, err := cmd.Output()
	if err != nil {
    	if *setup.LogLVL <= 2 {
			slog.Warn(fmt.Sprintf("%v", err))
		}
	}
	os.Exit(1)
}