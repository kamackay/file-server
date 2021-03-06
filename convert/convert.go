package convert

import (
	"github.com/google/uuid"
	"github.com/kamackay/goffmpeg/models"
	"github.com/kamackay/goffmpeg/transcoder"
	"github.com/sirupsen/logrus"
	"gitlab.com/kamackay/filer/files"
	"time"
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
	StartTime  int64     `json:"startTime"`
	InputFile  string    `json:"inputFile"`
	Duration   int64     `json:"duration"`
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

func (this *Converter) Convert(input string, request Request) uuid.UUID {
	jobId := uuid.New()
	this.jobs[jobId] = Job{
		Id:         jobId,
		StartTime:  time.Now().UnixNano(),
		Duration:   0,
		InputFile:  input,
		OutputFile: request.OutputFile,
		Status:     InProgress,
	}
	go func() {
		trans := new(transcoder.Transcoder)

		err := trans.Initialize(input, request.OutputFile)
		if err != nil {
			this.updateJob(jobId, 100, Failed, err)
		}

		trans.MediaFile().SetPreset(request.Preset)
		trans.MediaFile().SetThreads(1)
		trans.MediaFile().SetCRF(request.CRF)
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
			err = files.WriteMetaFileFor(request.OutputFile)
			if err != nil {
				this.log.Warnf("Error Writing Metadata for Converted File: %s", err)
			}
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
	job.Duration = time.Now().UnixNano() - job.StartTime
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
	Preset     string `json:"preset"`
	CRF        uint32 `json:"crf"`
}
