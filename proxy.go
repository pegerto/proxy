package main

import (
	"./http_proxy"
	"flag"
)

func main() {
	httpPort := flag.Int("httpPort", 8080, "Bind port for http proxy service")
	flag.Parse()

	httpProxy := http_proxy.HttpProxy{*httpPort}
	httpProxy.ListenAndServe()

}
