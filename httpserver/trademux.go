package httpserver

import (
	"net/http"
	"strings"

	"github.com/lworkltd/kits/service/httpsrv"
	"github.com/sirupsen/logrus"
)

type tradeMux struct {
	wrapper *httpsrv.Wrapper
	addr    string
}

func (mux *tradeMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/v1/trade/stream") {
		handleWs(w, req)
		return
	}
	mux.wrapper.ServeHTTP(w, req)
}

func (mux *tradeMux) run() (err error) {
	defer func() {
		if err == nil {
			return
		}
		logrus.Errorf("ListenAndServe failed,%v", err)
	}()

	log.WithField("on", mux.addr).Infof("Serve HTTP")
	err = http.ListenAndServe(mux.addr, mux)

	return
}
