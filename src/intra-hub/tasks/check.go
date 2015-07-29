package tasks

import (
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

func (dc *DatabaseCheck) isConnected() error {
	cmd := exec.Command("pidof", "mysql")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	if string(out) != "" {
		return nil
	}
	return fmt.Errorf("Can't connect to database")
}

func (dc *DatabaseCheck) Check() error {
	err := dc.isConnected()
	return err
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
