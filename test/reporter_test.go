package test

import (
	"github.com/latifrons/goerrorcollector/report"
	"testing"
	"time"
)

func Test(t *testing.T) {
	report.Start("amqp://guest:guest@localhost:5672/", "reporter", "reporter")
	report.Report("test", 1, "test", "test")
	time.Sleep(1 * time.Second)
}
