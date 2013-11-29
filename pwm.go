package bbio

import (
	"fmt"
	"os"
)

type pwmExport struct {
	periodFile   *os.File
	dutyFile     *os.File
	polarityFile *os.File
	duty         uint
	periodNs     uint
}

var (
	pwmInitialized bool
	ocpDir         string
	exportedPWMs   = make(map[string]*pwmExport)
)

func initializePWM() error {
	if pwmInitialized {
		return nil
	}

	err := LoadDeviceTree("am33xx_pwm")
	if err != nil {
		return err
	}

	ocpDir, err = buildPath("/sys/devices", "ocp")
	return err
}

func PWMSetFrequency(key string, freq float32) error {
	if freq <= 0 {
		return fmt.Errorf("invalid requency: %f", freq)
	}

	pwm, ok := exportedPWMs[key]
	if !ok {
		return fmt.Errorf("PWM '%s' not found", key)
	}

	periodNs := uint(1e9 / freq)
	if periodNs != pwm.periodNs {
		_, err := fmt.Fprintf(pwm.periodFile, "%d", periodNs)
		if err != nil {
			return err
		}
		pwm.periodNs = periodNs
	}

	return nil
}

func PWMSetPolarity(key string, polarity int) error {
	pwm, ok := exportedPWMs[key]
	if !ok {
		return fmt.Errorf("PWM '%s' not found", key)
	}

	_, err := fmt.Fprintf(pwm.polarityFile, "%d", polarity)
	return err
}

func PWMSetDutyCycle(key string, duty float32) error {
	if duty < 0 || duty > 100 {
		return fmt.Errorf("invalid duty cycle: %f", duty)
	}

	pwm, ok := exportedPWMs[key]
	if !ok {
		return fmt.Errorf("PWM '%s' not found", key)
	}

	pwm.duty = uint(float32(pwm.periodNs) * duty * 0.01)

	_, err := fmt.Fprintf(pwm.dutyFile, "%d", pwm.duty)
	return err
}

func PWMStart(key string, duty, freq float32, polarity int) error {
	err := initializePWM()
	if err != nil {
		return err
	}

	err = LoadDeviceTree("bone_pwm_" + key)
	if err != nil {
		return err
	}

	//finds and builds the pwmTestPath, as it can be variable...
	pwmTestPath, err := buildPath(ocpDir, "pwm_test_"+key)
	if err != nil {
		return err
	}

	//create the path for the period and duty
	periodPath := pwmTestPath + "/period"
	dutyPath := pwmTestPath + "/duty"
	polarityPath := pwmTestPath + "/polarity"

	periodFile, err := os.OpenFile(periodPath, os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	dutyFile, err := os.OpenFile(dutyPath, os.O_RDWR, 0660)
	if err != nil {
		periodFile.Close()
		return err
	}
	polarityFile, err := os.OpenFile(polarityPath, os.O_RDWR, 0660)
	if err != nil {
		periodFile.Close()
		dutyFile.Close()
		return err
	}

	exportedPWMs[key] = &pwmExport{
		periodFile:   periodFile,
		dutyFile:     dutyFile,
		polarityFile: polarityFile,
	}

	return nil
}

func PWMDisable(key string) {
	UnloadDeviceTree("bone_pwm_" + key)
	if pwm, ok := exportedPWMs[key]; ok {
		pwm.periodFile.Close()
		pwm.dutyFile.Close()
		pwm.polarityFile.Close()
		delete(exportedPWMs, key)
	}
}

func PWMCleanup() {
	for key := range exportedPWMs {
		PWMDisable(key)
	}
}
