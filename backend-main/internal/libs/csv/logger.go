package csv

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
)

type csvLogger struct {
	out io.Writer
}

var logger csvLogger

type CsvData struct {
	UserId          int64
	UserName        string
	SessionStart    time.Time
	SessionDuration time.Duration
}

func Setup(filepath string) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(errors.Wrap(err, "open csv log file"))
	}
	logger = csvLogger{out: file}
}

func Write(data CsvData) {
	startDateTime := data.SessionStart.Format(time.DateTime)
	duration := data.SessionDuration.String()

	s := fmt.Sprintf("%d;%s;%s;%s", data.UserId, data.UserName, startDateTime, duration)
	io.WriteString(logger.out, s)
}
