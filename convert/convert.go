package convert

import (
	"github.com/xfrr/goffmpeg/models"
	"github.com/xfrr/goffmpeg/transcoder"
)

type Converter struct {
}

func New() *Converter {
	return &Converter{}
}

func (this *Converter) Convert(input string,
	output string,
	progressUpdate func(progress models.Progress)) error {

	trans := new(transcoder.Transcoder)

	err := trans.Initialize(input, output)
	if err != nil {
		return err
	}

	done := trans.Run(true)

	progress := trans.Output()

	for msg := range progress {
		progressUpdate(msg)
	}

	err = <-done
	if err != nil {
		return err
	} else {
		return nil
	}
}

type Request struct {
	InputFile  string `json:"input"`
	OutputFile string `json:"output"`
}
