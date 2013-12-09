package bbio

import (
	"fmt"
	"os"
)

var (
	adcInitialized bool
	adcPrefixDir   string
)

type AIn string

const (
	AIN0 AIn = "AIN0"
	AIN1 AIn = "AIN1"
	AIN2 AIn = "AIN2"
	AIN3 AIn = "AIN3"
	AIN4 AIn = "AIN4"
	AIN5 AIn = "AIN5"
	AIN6 AIn = "AIN6"
)

func AInByPin(pinKey string) (AIn, bool) {
	pin, ok := PinByKey(pinKey)
	if !ok {
		return "", false
	}
	return AIn(fmt.Sprintf("AIN%d", pin.AIn)), true
}

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
	ain  AIn
	file *os.File
}

func NewADC(ain AIn) (*ADC, error) {
	if !adcInitialized {
		err := adcInit()
		if err != nil {
			return nil, err
		}
	}

	filename := adcPrefixDir + string(ain)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &ADC{ain, file}, nil
}

func (adc *ADC) AIn() AIn {
	return adc.ain
}

func (adc *ADC) ReadRaw() (value float32) {
	adc.file.Seek(0, os.SEEK_SET)
	fmt.Fscan(adc.file, &value)
	return value
}

func (adc *ADC) ReadValue() (value float32) {
	return adc.ReadRaw() / 1800.0
}

func (adc *ADC) Close() error {
	return adc.file.Close()
}

func CleanupADC() error {
	adcInitialized = false
	return UnloadDeviceTree("cape-bone-iio")
}
