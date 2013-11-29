package adc

import (
	"fmt"
	"os"

	"github.com/ungerik/go-bbio"
)

var (
	initialized bool
	prefixDir   string
)

func Setup() error {
	if initialized {
		return nil
	}

	err := bbio.LoadDeviceTree("cape-bone-iio")
	if err != nil {
		return err
	}

	ocpDir, _ := bbio.BuildPath("/sys/devices", "ocp")
	prefixDir, _ = bbio.BuildPath(ocpDir, "helper")
	prefixDir += "/AIN"
	initialized = true

	return nil
}

func ReadValue(ain uint) (value float32) {
	filename := fmt.Sprintf("%s%d", prefixDir, ain)
	file, _ := os.Open(filename)
	fmt.Fscan(file, &value)
	file.Close()
	return value
}

func Cleanup() error {
	initialized = false
	return bbio.UnloadDeviceTree("cape-bone-iio")
}
