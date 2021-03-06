package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/youngsterxyf/memcached-ui/config"
	"github.com/youngsterxyf/memcached-ui/controller"
)

const (
	VERSION = "0.2.0"
	APPNAME = "memcached-ui"
)

var (
	showv  bool
	listen string
	conf   string
	debug  bool
	// Git SHA Value will be set during build
	GitSHA    = "Not provided (use ./build instead of go build)"
	BuildTime = "Not provided (use ./build instead of go build)"
)

func init() {
	flag.BoolVar(&showv, "v", false, "show version of "+APPNAME)
	flag.StringVar(&listen, "l", ":8080", "memcached-ui server addr")
	flag.StringVar(&conf, "c", "config.yaml", "memcached-ui conf file")
	flag.BoolVar(&debug, "debug", false, "run in debug mode")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -c=config.yaml -l=:8080 -debug\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func showVersion() {
	fmt.Printf("%s v%s\n", APPNAME, VERSION)
	fmt.Printf("%10s : %s\n", "Built by", runtime.Version())
	fmt.Printf("%10s : %s\n", "Built at", BuildTime)
	fmt.Printf("%10s : %s\n", "Git SHA", GitSHA)
}

func appConfigMiddleware(conf config.AppConfigStruct) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("app_conf", conf)
		c.Next()
	}
}

func main() {
	flag.Parse()
	if showv {
		showVersion()
		return
	}

	appConfig, err := config.LoadAppConfig(conf)
	if err != nil {
		fmt.Printf("config load err: %s\n", err)
		os.Exit(-1)
	}

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Static("/assets", "./ui/assets")
	r.LoadHTMLGlob("ui/templates/*")
	r.Use(appConfigMiddleware(appConfig))

	if appConfig.Basic_auth.On {
		r.Use(gin.BasicAuth(gin.Accounts{appConfig.Basic_auth.Username: appConfig.Basic_auth.Password}))
	}

	r.GET("/", controller.Home)
	r.GET("/node", controller.Node)
	r.GET("/cluster", controller.Cluster)
	r.POST("/do", controller.Do)

	r.Run(listen)
}
