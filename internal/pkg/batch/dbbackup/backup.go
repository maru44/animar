package dbbackup

import (
	"animar/v1/configs"
	"fmt"
	"io/ioutil"
	"os/exec"
)

const backupFileName = "../../seed/backup/backup_main.sql"

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
	if err = ioutil.WriteFile(backupFileName, bytes, 0600); err != nil {
		fmt.Print(err)
	}
}

// func Bad() {
// 	cmd := exec.Command(
// 		"mysqldump", "--single-transaction", "--skip-lock-tables",
// 		fmt.Sprintf("-u%s", configs.MysqlUser),
// 		fmt.Sprintf("-p%s", configs.MysqlPassword),
// 		fmt.Sprintf("%s", configs.MysqlDataBase),
// 		fmt.Sprintf("> %s", backupFileName),
// 	)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	if err := cmd.Run(); err != nil {
// 		log.Print(err)
// 	}
// }

// func Sample() {
// 	cmd := exec.Command("pwd")
// 	if result, err := cmd.Output(); err != nil {
// 		log.Print(err)
// 	} else {
// 		log.Print(string(result))
// 	}
// }
