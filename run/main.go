package main

import (
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	log             = logging.MustGetLogger("netgo")
	httpListen      = flag.String("http", ":80", "host:port to listen on")
	quit            = make(chan bool)
	sleepForQuit    = 2
	saveRemoteTo, _ = os.Getwd()
)

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
// Helpers

func ProcessExecCmd(Cmd string) error {
	cmd := exec.Command("cmd", "/c", "start", Cmd)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//err := cmd.Run()
	err := cmd.Start()
	if nil != err {
		log.Debug("[ProcessExecCmd] [%+v]", err)
		return err
	}
	return nil
}
func ProcessExecFile(Url string) error {
	name_idx := strings.LastIndex(Url, "/")
	if -1 == name_idx {
		log.Debug("[ProcessExecFile] [No '/' in Url]")
		return errors.New("Invalid Filename")
	}

	filename := filepath.Join(saveRemoteTo, Url[name_idx+1:])
	log.Debug("[ProcessExecFile] [%v]", filename)

	if err := SaveFileFromUri(filename, Url); nil != err {
		log.Debug("[ProcessExecFile] [%+v]", err)
		return err
	}

	return ProcessExecCmd(filename)
}

func SaveFileFromUri(out_path string, in_path string) error {
	log.Debug("[SaveFileFromUri] [%+v] [%+v]", out_path, in_path)

	out, err := os.Create(out_path)
	if nil != err {
		log.Debug("[File not created: %+v]", err)
		return err
	}
	defer out.Close()

	resp, err := http.Get(in_path)
	if nil != err {
		log.Debug("[Path not found: %+v]", err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if nil != err {
		log.Debug("[Copy failed: %+v]", err)
		return err
	}

	return nil
}

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
// Route Handlers

func HandleShutdown(c *gin.Context) {
	log.Critical("[HandleShutdown]")

	c.Writer.Write([]byte("Shutting down..."))
	c.Writer.Flush()
	c.Request.Close = true
}
func HandleNetget(c *gin.Context) {

	in := struct {
		Url      string `form:"URL"`
		Key      string `form:"Key"`
		Function string `form:"_function"`
	}{}

	success := c.Bind(&in)
	if !success {
		c.XML(500, gin.H{"Result": "Incorrect Parameters"})
		return
	}

	log.Critical("[HandleNetget] [%+v]", in)

	if "ExecProtocol" == in.Function {
		err := ProcessExecCmd(in.Url)

		switch err {
		case nil:
			c.XML(200, gin.H{"Result": "Success"})
		default:
			c.XML(500, gin.H{"Result": "Fail", "Err": err.Error()})
		}
		return
	}
	if "ExecRemoteFile" == in.Function {
		err := ProcessExecFile(in.Url)

		switch err {
		case nil:
			c.XML(200, gin.H{"Result": "Success"})
		default:
			c.XML(500, gin.H{"Result": "Fail", "Err": err.Error()})
		}
		return
	}

	c.XML(404, gin.H{"Result": "Not Found"})
}

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
// Middleware

func ShutdownCb() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wait until all filters/handlers are done
		c.Next()

		// Notify quit
		close(quit)
	}
}

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
// Create Router Paths

func CreateHttpRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Define default route
	r.GET("/", HandleNetget)

	// Define shutdown route
	shutdown := r.Group("/shutdown")
	shutdown.Use(ShutdownCb())
	shutdown.GET("/", HandleShutdown)

	return r
}

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
// Main

func main() {
	log.Critical("[main start]")
	defer func() { log.Critical("[main end]") }()

	flag.Parse()

	InitLogging()

	// Listen for async Quit request
	go func() {
		select {
		case <-quit:
			log.Critical("[Quit Gracefully]")
			os.Exit(1)
		}
		time.Sleep(time.Duration(sleepForQuit))
	}()

	router := CreateHttpRouter()

	log.Critical("[Starting Service] [%s]", *httpListen)
	router.Run(*httpListen)
	//log.Fatal(http.ListenAndServe(*httpListen, router))
}

func InitLogging() {
	// Customize the output format
	logging.SetFormatter(logging.MustStringFormatter("▶ %{level:.1s} %{message}"))

	// Setup one stdout and one syslog backend.
	console_log := logging.NewLogBackend(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
	console_log.Color = false

	//var file_log_filter = syslog.LOG_DEBUG|syslog.LOG_LOCAL0|syslog.LOG_CRIT
	//log.Info("[%d]", file_log_filter)
	//file_log := logging.NewSyslogBackend("")
	//file_log.Color = false

	// Combine them both into one logging backend.
	logging.SetBackend(console_log)

	logging.SetLevel(logging.INFO, "netgo")
}
