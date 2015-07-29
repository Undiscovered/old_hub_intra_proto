package tasks

import (
	"fmt"
	"os/exec"

	"github.com/astaxie/beego/toolbox"
	"github.com/calmh/du"
	"github.com/pyk/byten"
	"os"
)

const (
	databaseCheck = "Database check"
	sizeCheck     = "Size check"
)

type DatabaseCheck struct {
}

func (dc *DatabaseCheck) isConnected() error {
	cmd := exec.Command("pidof", "mysql")
	out, _ := cmd.Output()
	if string(out) != "" {
		return nil
	}
	cmd = exec.Command("pidof", "mysqld")
	out, _ = cmd.Output()
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
	info, err := du.Get(os.Getenv("$HOME"))
	if err != nil {
		return err
	}
	return fmt.Errorf("Disk Total: %s - Disk Usage: %s - Disk Remaining: %s",
		byten.Size(info.TotalBytes), byten.Size(info.TotalBytes-info.FreeBytes), byten.Size(info.FreeBytes))
}

func init() {
	toolbox.AddHealthCheck(sizeCheck, &SizeCheck{})
	toolbox.AddHealthCheck(databaseCheck, &DatabaseCheck{})
}
