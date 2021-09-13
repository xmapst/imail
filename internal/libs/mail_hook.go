package libs

import (
	"errors"
	"fmt"
	"github.com/midoks/imail/internal/config"
	"os"
	"os/exec"
)

func ExecPython(scriptName string, id int64) (string, error) {
	hookEnable, _ := config.GetBool("hook.enable", false)
	if !hookEnable {
		return "", errors.New("config is disable!")
	}

	cpath, _ := os.Getwd()
	fileName := fmt.Sprintf("%s/hook/%s", cpath, scriptName)
	_, b := IsExists(fileName)
	// fmt.Println(fileName, b)
	if !b {
		return "", errors.New("file is not exist!")
	}

	cmd := exec.Command("python", fileName, string(id))
	out, err := cmd.CombinedOutput()
	return string(out), err
}