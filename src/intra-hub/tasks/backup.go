package tasks

import "github.com/astaxie/beego/toolbox"

const (
	backupTaskName = "Backup"
)

func backup() error {
	return nil
}

func init() {
	crawler := toolbox.NewTask(blowFishCrawlerTaskName, "0 0 8 * * 1", backup)
	toolbox.AddTask(backupTaskName, crawler)
	toolbox.StartTask()
}
