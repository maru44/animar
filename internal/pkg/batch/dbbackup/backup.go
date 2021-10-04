package dbbackup

import (
	"animar/v1/configs"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func BackupMainDatabase() {

	cmd := exec.Command(
		"mysqldump", "--single-transaction", "--skip-lock-tables",
		fmt.Sprintf("-u%s", configs.MysqlUser),
		fmt.Sprintf("-p%s", configs.MysqlPassword),
		fmt.Sprintf("%s", configs.MysqlDataBase),
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Print(err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Print(err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Print(err)
	}
	if err = ioutil.WriteFile(configs.MysqlBackupFile, bytes, 0600); err != nil {
		fmt.Print(err)
	}
}
