package webserver

import (
    "github.com/gin-gonic/gin"
    "github.com/losenmann/iptv-toolkit/udpxy"
    "net/http"
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
    router.StaticFS(p, http.Dir(d))
}

func Run(p int) {
    router.Run(":" + strconv.Itoa(p))
}
