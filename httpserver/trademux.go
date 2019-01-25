package httpserver

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type tradeMux struct {
	ginEngine *gin.Engine
	addr      string
}

func (mux *tradeMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/v1/trade/stream") {
		handleWs(w, req)
		return
	}

	mux.ginEngine.ServeHTTP(w, req)
}

func (mux *tradeMux) run() (err error) {
	defer func() {
		if err == nil {
			return
		}
		logrus.Errorf("ListenAndServe failed,%v", err)
	}()
	logrus.Infof("Listening and serving HTTP on %s\n", mux.addr)
	err = http.ListenAndServe(mux.addr, mux)
	return
}
