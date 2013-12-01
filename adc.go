package bbio

import (
	"fmt"
	"os"
)

var (
	adcInitialized bool
	adcPrefixDir   string
)

func adcInit() error {
	err := LoadDeviceTree("cape-bone-iio")
	if err != nil {
		return err
	}

	ocpDir, _ := BuildPath("/sys/devices", "ocp")
	adcPrefixDir, _ = BuildPath(ocpDir, "helper")
	adcPrefixDir += "/AIN"
	adcInitialized = true

	return nil
}

type ADC struct {
	ain  int
	file *os.File
}

func NewADC(ain int) (*ADC, error) {
	if !adcInitialized {
		err := adcInit()
		if err != nil {
			return nil, err
		}
	}

	filename := fmt.Sprintf("%s%d", adcPrefixDir, ain)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &ADC{ain, file}, nil
}

func (self *ADC) AIn() int {
	return self.ain
}

func (self *ADC) ReadValue() (value float32) {
	self.file.Seek(0, os.SEEK_SET)
	fmt.Fscan(self.file, &value)
	return value
}

func (self *ADC) Close() error {
	return self.file.Close()
}

func CleanupADC() error {
	adcInitialized = false
	return UnloadDeviceTree("cape-bone-iio")
}
