package logger

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Configure() error {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{
		// FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fileName := f.File[strings.LastIndex(f.File, "/")+1:] + ":" + strconv.Itoa(f.Line)
			fnName := f.Function[strings.LastIndex(f.Function, ".")+1:]
			return fnName, fileName
		},
	})

	if viper.GetBool("log.debug") {
		log.SetLevel(log.DebugLevel)
	}

	// check if we need to store log to file
	if viper.GetBool("log.file.enable") {
		if viper.GetString("log.file.path") == "" {
			fmt.Println("You enable the file logging, but not define the log path")
			os.Exit(1)
		}

		err := createDirectoryByFile(viper.GetString("log.file.path"))
		if err != nil {
			return err
		}

		file, err := os.OpenFile(viper.GetString("log.file.path"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		log.SetOutput(file)
	}

	return nil
}

func createDirectoryByFile(path string) error {
	dirs := strings.Split(path, "/")

	// If only file name or the string is empty end the process
	if len(dirs) <= 1 {
		return nil
	}

	dirs = dirs[:len(dirs)-1]
	dir := strings.Join(dirs, "/")

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
