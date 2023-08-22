package klog

import (
	"flag"
	"io"
	"istomyang.github.com/like-iam/log"
	"k8s.io/klog/v2"
)

//InfoLog:    "INFO",
//WarningLog: "WARNING",
//ErrorLog:   "ERROR",
//FatalLog:   "FATAL",

func Init(log *log.Logger) {
	// https://pkg.go.dev/k8s.io/klog/v2#readme-how-to-use-klog

	fs := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(fs)
	defer klog.Flush()

	klog.SetOutputBySeverity("INFO", &infoLogger{log: log})
	klog.SetOutputBySeverity("WARNING", &warnLogger{log: log})
	klog.SetOutputBySeverity("ERROR", &errorLogger{log: log})
	klog.SetOutputBySeverity("FATAL", &fatalLogger{log: log})

	// Here you can refer to source klog.init() function.
	_ = fs.Set("skip_headers", "true")
	_ = fs.Set("logtostderr", "false")
}

type infoLogger struct {
	log *log.Logger
}

func (i *infoLogger) Write(p []byte) (n int, err error) {
	// In klog src, p[len(p)] === '\n', here should remove it.
	i.log.Info(string(p[:len(p)-1]))

	// interface Writer says: Must return len(p) if no error.
	return len(p), nil
}

var _ io.Writer = &infoLogger{}

type warnLogger struct {
	log *log.Logger
}

func (w *warnLogger) Write(p []byte) (n int, err error) {
	w.log.Warn(string(p[:len(p)-1]))
	return len(p), nil
}

var _ io.Writer = &warnLogger{}

type errorLogger struct {
	log *log.Logger
}

func (e *errorLogger) Write(p []byte) (n int, err error) {
	e.log.Error(string(p[:len(p)-1]))
	return len(p), nil
}

var _ io.Writer = &errorLogger{}

type fatalLogger struct {
	log *log.Logger
}

func (f *fatalLogger) Write(p []byte) (n int, err error) {
	f.log.Fatal(string(p[:len(p)-1]))
	return len(p), nil
}

var _ io.Writer = &fatalLogger{}
