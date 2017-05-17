package qframe_collector_file

import (
	"log"
	"os"

	"github.com/hpcloud/tail"
	"github.com/qnib/qframe-types"
	"github.com/zpatrick/go-config"
)

const (
	version = "0.0.2"
	pluginTyp = "collector"
	pluginPkg = "file"
)

type Plugin struct {
	QChan qtypes.QChan
	Cfg config.Config
	Name string
}

func NewPlugin(qChan qtypes.QChan, cfg *config.Config, name string) Plugin {
	return Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
	}
}

func (p *Plugin) Run() {
	log.Printf("[II] Start file collector v%s", version)
	fPath, err := p.Cfg.String("collector.file.path")
	if err != nil {
		log.Println("[EE] No file path for collector.file.path set")
		return
	}
	create, _ := p.Cfg.BoolOr("collector.file.create", false)
	if _, err := os.Stat(fPath); os.IsNotExist(err) && create {
		log.Printf("[DD] Create file: %s", fPath)
		f, _ := os.Create(fPath)
		f.Close()
	}
	fileReopen, err := p.Cfg.BoolOr("collector.file.reopen", true)
	t, err := tail.TailFile(fPath, tail.Config{Follow: true, ReOpen: fileReopen})
	if err != nil {
		log.Printf("[WW] File collector failed to open %s: %s", fPath, err)
	}
	for line := range t.Lines {
		qm := qtypes.NewQMsg("collector", p.Name)
		qm.Msg = line.Text
		p.QChan.Data.Send(qm)
	}
}
