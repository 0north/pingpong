package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

var (
	PING_URL   = os.Getenv("PING_URL")
	DD_SERVICE = os.Getenv("DD_SERVICE")
	DD_ENV     = os.Getenv("DD_ENV")
)

func ping() {
	for {
		resp, err := http.Get(PING_URL)
		if err != nil {
			logrus.Info("cannot call ping: %s\n", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		logrus.Info("response from %v: %s\n", PING_URL, string(body))
		time.Sleep(5 * time.Second)
	}
}

func pong() {
	r := gin.Default()
	// Add Datadog tracing middleware
	r.Use(gintrace.Middleware(DD_SERVICE, gintrace.WithResourceNamer(func(ctx *gin.Context) string {
		return ctx.Request.Method + " " + ctx.Request.URL.Path
	})))

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(&DDContextLogHook{})
	logrus.SetOutput(os.Stdout)

	err := profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			// profiler.BlockProfile,
			// profiler.MutexProfile,
			// profiler.GoroutineProfile,
		),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	defer profiler.Stop()

	r.GET("/ping", func(c *gin.Context) {
		logrus.WithContext(c.Request.Context()).Info("served pong")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	tracer.Start()
	defer tracer.Stop()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func main() {
	fmt.Println("PING_URL: " + PING_URL)
	if PING_URL != "" {
		ping()
	} else {
		pong()
	}
}

type DDContextLogHook struct{}

// Levels implements logrus.Hook interface, this hook applies to all defined levels
func (d *DDContextLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implements logrus.Hook interface, attaches trace and span details found in entry context
func (d *DDContextLogHook) Fire(e *logrus.Entry) error {
	span, _ := tracer.SpanFromContext(e.Context)
	e.Data["dd.trace_id"] = span.Context().TraceID()
	e.Data["dd.span_id"] = span.Context().SpanID()
	return nil
}
