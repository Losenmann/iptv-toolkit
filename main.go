package main

import (
    "fmt"
    "github.com/losenmann/iptv-toolkit/scheduler"
    "github.com/losenmann/iptv-toolkit/setup"
    "github.com/losenmann/iptv-toolkit/webserver"
    "log/slog"
    "os"
)

func init() {
    setup.Initgo()
}

func main() {
    if err := os.MkdirAll(*setup.WebDir, 0744); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
        os.Exit(1)
    }

    scheduler.Task()
    if *setup.Schedule {
        scheduler.Main(*setup.Crontab)
    }
    if *setup.Udpxy {
        webserver.Udpxy(*setup.WebPathUdpxy)
    }
    if *setup.Files {
        webserver.Files(*setup.WebPathFiles, *setup.WebDir)
    }
    webserver.Run(*setup.WebPort)
}
