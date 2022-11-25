//***************************//
// Package logger
// Date Of Creation: 11/01/2021
// Authour Name: Moloy Mondal
// History: N/A
// Synopsis: This file is used to set logger configuration
// Functions: init,
// Global Variable: Log
// Version: 1.0.0
//***************************//
package logger

import (
	"flag"
	"log"
	"os"
	"github.com/natefinch/lumberjack"
	_ "github.com/natefinch/lumberjack"
)

var (
	Log *log.Logger
)

func init() {
	// set location of log file
	var logpath = "log/logfile.log"

	flag.Parse()
	var file, err1 = os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err1 != nil {
		panic(err1)
	}
	//Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	//Log.Println("LogFile : " + logpath)
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Log.SetOutput(&lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    25, // megabytes after which new file is created
		MaxBackups: 3,  // number of backups
		MaxAge:     3,  //days
	})
}
