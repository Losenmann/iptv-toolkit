package main

import (
    "log/slog"
    "os"
    "fmt"
    "iptv-toolkit/main/setup"
    "iptv-toolkit/main/webserver"
    "iptv-toolkit/main/scheduler"
    "iptv-toolkit/main/udpxy"
)

func init() {
    setup.Initgo()
}

func main() {
    if err := os.MkdirAll(*setup.WebPath, 0744); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        os.Exit(1)
    }

    scheduler.Task()
    if *setup.Schedule == true {
        scheduler.Main(*setup.Crontab)
    }
    if *setup.Udpxy == true {
        go udpxy.UdpxyExt()
    }
    if *setup.Web == true {
        webserver.Main(*setup.WebPort, *setup.WebPath)
    }
}

