package qframe_collector_file

import (
	"C"
	"log"

	"github.com/hpcloud/tail"
	"github.com/qnib/qframe/types"
	"github.com/zpatrick/go-config"
)

const (
	version = "0.0.1"
)

type Plugin struct {
	QChan qtypes.QChan
	Cfg config.Config
}

func NewPlugin(qChan qtypes.QChan, cfg config.Config) Plugin {
	return Plugin{
		QChan: qChan,
		Cfg: cfg,
	}
}

func (p *Plugin) Run() {
	log.Printf("[II] Start file collector v%s", version)
	fPath, err := p.Cfg.String("collector.file.path")
	if err != nil {
		log.Println("[EE] No file path for collector.file.path set")
		return
	}
	fileReopen, err := p.Cfg.BoolOr("collector.file.reopen", true)
	t, err := tail.TailFile(fPath, tail.Config{Follow: true, ReOpen: fileReopen})
	if err != nil {
		log.Printf("[WW] File collector failed to open %s: %s", fPath, err)
	}
	for line := range t.Lines {
		qm := qtypes.NewQMsg("collector", "fPath")
		qm.Msg = line.Text
		p.QChan.Data.Send(qm)
	}
}
