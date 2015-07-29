package tasks

import (
	"os"
	"os/exec"
	"time"

	"intra-hub/confperso"

	"github.com/astaxie/beego/toolbox"
)

const (
	backupTaskName = "Backup"
)

func backup() error {
	fileName := time.Now().Format(time.RFC3339) + ".sql"
	outputDirectoryName := os.Getenv("HOME") + "/backup"
	os.Mkdir(outputDirectoryName, 0775)
	outputFileName := outputDirectoryName + "/" + fileName
	// generates something like mysqldump -u root --p="password" intra_hub > /Users/Vincent/backup/2015-07-29T13:58:36+02:00.sql
	cmd := "mysqldump -u " + confperso.Username + " --password=\"" + confperso.Password + "\" " + confperso.DatabaseName + " > " + outputFileName
	_, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return err
}

func init() {
	crawler := toolbox.NewTask(blowFishCrawlerTaskName, "0 0 8 * * 1", backup)
	toolbox.AddTask(backupTaskName, crawler)
	toolbox.StartTask()
}
