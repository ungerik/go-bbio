package bbio

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ungerik/go-quick"
)

// var gpio_mode int
// var gpio_direction [120]int
// var pwm_pins [120]int
var (
	ctrlDir string
)

func buildPath(partialPath, prefix string) (path string, found bool) {
	dirFiles, err := ioutil.ReadDir(partialPath)
	if err != nil {
		return "", false
	}
	for _, file := range dirFiles {
		if file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
			return file.Name(), true
		}
	}
	return "", false
}

func LoadDeviceTree(name string) error {
	ctrlDir, _ = buildPath("/sys/devices", "bone_capemgr")
	slots := ctrlDir + "/slots"

	data, err := quick.FileGetString(slots)
	if err != nil {
		return err
	}

	if strings.Contains(data, name) {
		return nil
	}

	err = quick.FileSetString(slots, name)
	if err == nil {
		time.Sleep(time.Millisecond * 200)
	}
	return err
}

func UnloadDeviceTree(name string) error {
	slots := ctrlDir + "/slots"

	file, err := os.OpenFile(slots, os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	for err != nil {
		if strings.Contains(line, name) {
			line = line[:strings.IndexRune(line, ':')]
			line = strings.TrimSpace(line)
			_, err = file.WriteString("-" + line)
			return err
		}
		line, err = reader.ReadString('\n')
	}
	if err != io.EOF {
		return err
	}
	return nil
}
