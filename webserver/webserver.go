package webserver

import (
    "github.com/gin-gonic/gin"
    "github.com/losenmann/iptv-toolkit/udpxy"
    "strconv"
)

var router *gin.Engine

func init() {
    router = gin.Default()
}

func Udpxy(p string) {
    router.GET(p, udpxy.Udpxy)
}

func Files(p, d string) {
    router.Static(p, d)
}

func Run(p int) {
    router.Run(":" + strconv.Itoa(p))
}
