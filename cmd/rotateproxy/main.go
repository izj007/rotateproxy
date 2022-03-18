package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/akkuman/rotateproxy"
)

var (
	baseCfg     rotateproxy.BaseConfig
	listenaddr  string
	email       string
	token       string
	rule        string
	pageCount   int
	proxy       string
	portPattern = regexp.MustCompile(`^\d+$`)
)

func init() {
	flag.StringVar(&baseCfg.ListenAddr, "l", ":8899", "listen address")
	flag.StringVar(&listenaddr, "lw", ":9000", "Http Server listen address")
	flag.StringVar(&baseCfg.Username, "user", "", "authentication username")
	flag.StringVar(&baseCfg.Password, "pass", "", "authentication password")
	flag.StringVar(&email, "email", "", "email address")
	flag.StringVar(&token, "token", "", "token")
	flag.StringVar(&proxy, "proxy", "", "proxy")
	flag.StringVar(&rule, "rule", "", `protocol=="socks5" && "Version:5 Method:No Authentication(0x00)" && after="2022-03-01" && country="CN"`)
	flag.IntVar(&baseCfg.IPRegionFlag, "region", 0, "0: all 1: cannot bypass gfw 2: bypass gfw")
	flag.IntVar(&baseCfg.SelectStrategy, "strategy", 99, "0: random, 1: Select the one with the shortest timeout, 2: Select the two with the shortest timeout, 3: Select the one without used...")
	flag.IntVar(&pageCount, "page", 5, "the page count you want to crawl")
	flag.Parse()
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	if !isFlagPassed("email") || !isFlagPassed("token") {
		flag.Usage()
		return
	}

	baseCfg.ListenAddr = strings.TrimSpace(baseCfg.ListenAddr)

	if portPattern.Match([]byte(baseCfg.ListenAddr)) {
		baseCfg.ListenAddr = ":" + baseCfg.ListenAddr
	}
	if rule == "" {
		nTime := time.Now()
		yesTime := nTime.AddDate(0, 0, -1)
		logDay := yesTime.Format("2006-01-02")
		rule = fmt.Sprintf(`protocol=="socks5" && "Version:5 Method:No Authentication(0x00)" && after="%s"`, logDay)
		fmt.Printf("[-]获取%s之后的代理", logDay)
	}

	rotateproxy.StartRunCrawler(token, email, rule, pageCount, proxy)
	rotateproxy.StartCheckProxyAlive()
	rotateproxy.StartWebServe(listenaddr)
	c := rotateproxy.NewRedirectClient(rotateproxy.WithConfig(&baseCfg))
	c.Serve()
	select {}
}
