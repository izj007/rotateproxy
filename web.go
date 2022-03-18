package rotateproxy

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func StartWebServe(listenaddr string) {
	http.HandleFunc("/getproxy", doGetProxy) //   设置访问路由
	go http.ListenAndServe(listenaddr, nil)
}

func doGetProxy(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key, err := GetRandomProxyUrl()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	key = strings.TrimPrefix(key, "socks5://")
	fmt.Fprintf(w, key)
}
