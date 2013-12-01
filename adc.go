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

func (adc *ADC) AIn() int {
	return adc.ain
}

func (adc *ADC) ReadValue() (value float32) {
	adc.file.Seek(0, os.SEEK_SET)
	fmt.Fscan(adc.file, &value)
	return value
}

func (adc *ADC) Close() error {
	return adc.file.Close()
}

func CleanupADC() error {
	adcInitialized = false
	return UnloadDeviceTree("cape-bone-iio")
}
