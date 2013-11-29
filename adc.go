package bbio

import (
	"fmt"
	"os"
)

var (
	adcInitialized bool
	adcPrefixDir   string
)

func ADCSetup() error {
	if adcInitialized {
		return nil
	}

	err := LoadDeviceTree("cape-bone-iio")
	if err != nil {
		return err
	}

	ocpDir, _ := buildPath("/sys/devices", "ocp")
	adcPrefixDir, _ = buildPath(ocpDir, "helper")
	adcPrefixDir += "/AIN"
	adcInitialized = true

	return nil
}

func ReadValue(ain int) (value float32) {
	filename := fmt.Sprintf("%s%d", adcPrefixDir, ain)
	file, _ := os.Open(filename)
	fmt.Fscan(file, &value)
	file.Close()
	return value
}

func ADCCleanup() error {
	adcInitialized = false
	return UnloadDeviceTree("cape-bone-iio")
}
