package tasks

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/astaxie/beego/toolbox"
	"github.com/canthefason/kite/systeminfo"
	"github.com/pyk/byten"
)

const (
	databaseCheck = "Database check"
	sizeCheck     = "Size check"
)

type DatabaseCheck struct {
}

func (dc *DatabaseCheck) isConnected() bool {
	cmd := exec.Command("pidof", "mysql")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	if string(out) != "" {
		return true
	}
	return false
}

func (dc *DatabaseCheck) Check() error {
	if dc.isConnected() {
		return nil
	} else {
		return errors.New("can't connect database")
	}
}

type SizeCheck struct {
}

func (dc *SizeCheck) Check() error {
	info, err := systeminfo.New()
	if err != nil {
		return err
	}
	return fmt.Errorf("Disk Total: %s - Disk Usage: %s - Disk Remaining: %s",
		byten.Size(int64(info.DiskTotal)), byten.Size(int64(info.DiskUsage)), byten.Size(int64(info.DiskTotal-info.DiskUsage)))
}

func init() {
	toolbox.AddHealthCheck(sizeCheck, &SizeCheck{})
	toolbox.AddHealthCheck(databaseCheck, &DatabaseCheck{})
}
