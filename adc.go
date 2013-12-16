package bbio

import (
	"fmt"
	"os"
)

var (
	adcInitialized bool
	adcPrefixDir   string
)

type AInName string

const (
	AIN0 AInName = "AIN0"
	AIN1 AInName = "AIN1"
	AIN2 AInName = "AIN2"
	AIN3 AInName = "AIN3"
	AIN4 AInName = "AIN4"
	AIN5 AInName = "AIN5"
	AIN6 AInName = "AIN6"
)

func AInNameByPin(pinKey string) (AInName, bool) {
	pin, ok := PinByKey(pinKey)
	if !ok {
		return "", false
	}
	return AInName(fmt.Sprintf("AIN%d", pin.AInNr)), true
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
	ain  AInName
	file *os.File
}

func NewADC(ain AInName) (*ADC, error) {
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

func (adc *ADC) AIn() AInName {
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
