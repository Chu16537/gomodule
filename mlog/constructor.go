package mlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Chu16537/module_master/mtime"
	"github.com/sirupsen/logrus"
)

type Config struct {
	LogFilePath  string
	LogFileName  string
	ElasticURL   string
	ElasticIndex string
}

type Log struct {
	handler *handler
	name    string
}
type handler struct {
	config      *Config
	file        *os.File
	t           time.Time
	currentDate string // 目前日期
}

var h *handler

func New(config *Config, t time.Time) error {
	opt := &logOpt{
		JSONFormatter: logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		},
	}

	logrus.SetFormatter(opt)
	logrus.SetLevel(logrus.DebugLevel)

	h = &handler{
		config:      config,
		t:           t,
		currentDate: mtime.GetTimeFormatTime(t, mtime.Format_YMD),
	}

	// 輸出到本地
	if config.LogFilePath != "" {
		fp := filepath.Join(config.LogFilePath, fmt.Sprintf("%s_%s.log", config.LogFileName, h.currentDate))

		// 確保目錄存在
		err := os.MkdirAll(filepath.Dir(fp), os.ModePerm)
		if err != nil {
			return err
		}

		h.file, err = os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			return err
		}

		logrus.SetOutput(h.file)
	} else {
		logrus.SetOutput(os.Stdout)
	}

	// 添加 Elasticsearch hook
	if config.ElasticURL != "" {
		esHook, err := newElasticsearchHook(config.ElasticURL, config.ElasticIndex)
		if err != nil {
			return err
		}

		logrus.AddHook(esHook)
	}

	return nil
}

func Get(name string) ILog {
	return &Log{
		handler: h,
		name:    name,
	}
}

func Done() {
	if h.file != nil {
		h.file.Close()
	}
}
