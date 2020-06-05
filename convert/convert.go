package convert

import (
	"fmt"
	"github.com/xfrr/goffmpeg/models"
	"github.com/xfrr/goffmpeg/transcoder"
)

type Converter struct {
	progressUpdate func(progress models.Progress)
}

func New(progressUpdate func(progress models.Progress)) *Converter {
	return &Converter{
		progressUpdate: progressUpdate,
	}
}

func (this *Converter) Convert(input string, output string) error {

	trans := new(transcoder.Transcoder)

	err := trans.Initialize(input, output)
	if err != nil {
		return err
	}

	done := trans.Run(true)

	progress := trans.Output()

	for msg := range progress {
		this.progressUpdate(msg)
		fmt.Println("Progress Update Received")
	}

	err = <-done
	if err != nil {
		return err
	} else {
		return nil
	}
}

type Request struct {
	OutputFile string `json:"output"`
}
