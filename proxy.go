package main

import ("io"
        "flag"
        "strconv"
        "net/http"
        "encoding/json"
        "time"
        "fmt"
      )

type handleConnection struct{}

type Transaction struct {
    RemoteAddr  string
    URL  string
    SC int32
}

func (*handleConnection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  tr := http.Transport{}
  tr.DisableKeepAlives = true
  resp , err := tr.RoundTrip(r)

  if err == nil {
    copyHeaders(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

    if bodyError := resp.Body.Close(); err != nil {
			http.Error(w, bodyError.Error(), 500)
      log(r, 500)
		}else{
      log(r, int32(resp.StatusCode))
    }
  }else {
    http.Error(w, err.Error(), 500)
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

func log(r *http.Request, sC int32){
  tr := Transaction{r.RemoteAddr, r.URL.String(), sC}
  line, _ := json.Marshal(tr)
  fmt.Println(string(line))
}


func main() {
    httpPort := flag.Int("httpPort", 8080, "Bind port for http proxy service")
    flag.Parse()
    s := &http.Server{
	     Addr:           ":"+ strconv.Itoa(*httpPort),
       Handler:        &handleConnection{},
	     ReadTimeout:    10 * time.Second,
	     WriteTimeout:   10 * time.Second,
	     MaxHeaderBytes: 1 << 20,
     }
     s.ListenAndServe()
}
