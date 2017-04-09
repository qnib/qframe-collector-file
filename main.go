package main

import (
	"C"
	"log"

	"github.com/hpcloud/tail"
	"github.com/qnib/qframe/types"
	"github.com/zpatrick/go-config"
)

const (
	version = "0.0.0.0"
)

func Run(qChan qtypes.QChan, cfg config.Config) {
	log.Printf("[II] Start file collector v%s", version)
	fPath, err := cfg.String("collector.file.path")
	if err != nil {
		log.Println("[EE] No file path for collector.file.path set")
		return
	}
	fileReopen, err := cfg.BoolOr("collector.file.reopen", true)
	t, err := tail.TailFile(fPath, tail.Config{Follow: true, ReOpen: fileReopen})
	if err != nil {
		log.Printf("[WW] File collector failed to open %s: %s", fPath, err)
	}
	for line := range t.Lines {
		qm := qtypes.NewQMsg("collector", "fPath")
		qm.Msg = line.Text
		qChan.Data.Send(qm)
	}
}
