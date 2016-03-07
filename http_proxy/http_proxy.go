package http_proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type handleConnection struct{}

type Transaction struct {
	RemoteAddr string
	Schema     string
	Host       string
	Path       string
	Query      string
	UserAgent  string
	SC         int32
}

func (*handleConnection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tr := http.Transport{}
	tr.DisableKeepAlives = true
	resp, err := tr.RoundTrip(r)

	if err == nil {
		copyHeaders(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

		if bodyError := resp.Body.Close(); err != nil {
			http.Error(w, bodyError.Error(), 500)
		} else {
			log(r, int32(resp.StatusCode))
		}
	}
}

func copyHeaders(dst, src http.Header) {
	for k, _ := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func log(r *http.Request, sC int32) {
	tr := Transaction{r.RemoteAddr,
		r.URL.Scheme,
		r.URL.Host,
		r.URL.Path,
		r.URL.Query().Encode(),
		r.UserAgent(),
		sC}

	line, _ := json.Marshal(tr)
	fmt.Println(string(line))
}

type HttpProxy struct {
	HttpPort int
}

func (httpProxy HttpProxy) ListenAndServe() {
	s := &http.Server{
		Addr:           ":" + strconv.Itoa(httpProxy.HttpPort),
		Handler:        &handleConnection{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
