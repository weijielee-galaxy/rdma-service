package main

import (
	"flag"
	"io"
	"log"
	"os"
	"rdma-service/cmd/acs"
	"rdma-service/cmd/kernelmodule"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// command line define
	logfile := flag.String("log", "./rdma-service.log", "logfile path")
	termi := flag.Bool("termi", false, "Print log to terminaind file")
	flag.Parse()

	//log setting
	var logOutput io.Writer
	if *termi {
		logOutput = io.MultiWriter(os.Stdout, &lumberjack.Logger{
			Filename:   *logfile,
			MaxSize:    10,
			MaxBackups: 10,
			MaxAge:     28,
			Compress:   true,
		})
	} else {
		logOutput = &lumberjack.Logger{
			Filename:   *logfile,
			MaxSize:    10,
			MaxBackups: 10,
			MaxAge:     28,
			Compress:   true}
	}

	log.SetOutput(logOutput)
	log.Printf("==========begin run the rdma service=============")

	ticker := time.NewTicker(time.Duration(10) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("============run the rdma service============")
		start0 := time.Now()
		acs.CheckACS()
		log.Printf("run the acs check elaspe: %v\n", time.Since(start0))

		start := time.Now()
		kernelmodule.CheckKernelMod()
		elaspe := time.Since(start)
		log.Printf("run the kernel module check elaspe:%v\n", elaspe)
		log.Printf("run the all check elaspe:%v\n", time.Since(start0))
	}
}
