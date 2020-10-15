package main

import ("KalbiProxy/pkg/proxy"
        "flag")

func main() {
	host := flag.String("host", "127.0.0.1", "IP Interface to bind proxy to")
	port := flag.Int("port", 5060, "Port to bind proxy to")
	flag.Parse()
	proxy := new(proxy.Proxy)
	proxy.Start(*host, *port)
}
