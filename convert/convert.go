package convert

import (
	"github.com/google/uuid"
	"github.com/kamackay/goffmpeg/models"
	"github.com/kamackay/goffmpeg/transcoder"
	"github.com/sirupsen/logrus"
)

const (
	InProgress int = iota
	Failed
	Done
)

type ProgressFunc = func(progress models.Progress, job Job)

type Converter struct {
	log            *logrus.Logger
	progressUpdate ProgressFunc
	jobs           map[uuid.UUID]Job
}

type Job struct {
	Id         uuid.UUID `json:"id"`
	InputFile  string    `json:"inputFile"`
	OutputFile string    `json:"outputFile"`
	Status     int       `json:"status"`
	Progress   float64   `json:"progress"`
	Error      error     `json:"error"`
}

func New(progressUpdate ProgressFunc) *Converter {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	return &Converter{
		log:            log,
		progressUpdate: progressUpdate,
		jobs:           map[uuid.UUID]Job{},
	}
}

func (this *Converter) Convert(input string, output string) uuid.UUID {
	jobId := uuid.New()
	this.jobs[jobId] = Job{
		Id:         jobId,
		InputFile:  input,
		OutputFile: output,
		Status:     InProgress,
	}
	go func() {
		trans := new(transcoder.Transcoder)

		err := trans.Initialize(input, output)
		if err != nil {
			this.updateJob(jobId, 100, Failed, err)
		}

		trans.MediaFile().SetPreset("veryfast")
		trans.MediaFile().SetThreads(2)
		//trans.MediaFile().SetBufferSize(200000)

		done := trans.Run(true)

		progress := trans.Output()

		for msg := range progress {
			this.progressUpdate(msg, this.MustGetJob(jobId))
			this.updateJob(jobId, msg.Progress, InProgress, nil)
		}

		err = <-done
		if err != nil {
			this.updateJob(jobId, 0, Failed, err)
		} else {
			this.updateJob(jobId, 100, Done, nil)
		}
	}()
	return jobId
}

func (this *Converter) GetJobStr(id string) *Job {
	return this.GetJob(uuid.MustParse(id))
}

func (this *Converter) GetJob(id uuid.UUID) *Job {
	if val, ok := this.jobs[id]; ok {
		return &val
	} else {
		return nil
	}
}

func (this *Converter) MustGetJob(id uuid.UUID) Job {
	return *this.GetJob(id)
}

func (this *Converter) updateJob(id uuid.UUID, progress float64, status int, err error) {
	job := this.jobs[id]
	if err != nil || status == Failed {
		this.log.Warnf("Job %s Has Failed %s", id.String(), err)
		job.Status = Failed
		job.Error = err
	} else {
		job.Status = status
		job.Error = nil
		job.Progress = progress
	}
	this.jobs[id] = job
}

type Request struct {
	OutputFile string `json:"output"`
}
