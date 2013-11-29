package pwm

import (
	"fmt"
	"os"

	"github.com/ungerik/go-bbio"
)

type export struct {
	periodFile   *os.File
	dutyFile     *os.File
	polarityFile *os.File
	duty         uint
	periodNs     uint
}

var (
	initialized bool
	ocpDir      string
	exported    = make(map[string]*export)
)

func initialize() error {
	if initialized {
		return nil
	}

	err := bbio.LoadDeviceTree("am33xx_pwm")
	if err != nil {
		return err
	}

	ocpDir, err = bbio.BuildPath("/sys/devices", "ocp")
	return err
}

func SetFrequency(key string, freq float32) error {
	if freq <= 0 {
		return fmt.Errorf("invalid requency: %f", freq)
	}

	pwm, ok := exported[key]
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

func SetPolarity(key string, polarity int) error {
	pwm, ok := exported[key]
	if !ok {
		return fmt.Errorf("PWM '%s' not found", key)
	}

	_, err := fmt.Fprintf(pwm.polarityFile, "%d", polarity)
	return err
}

func SetDutyCycle(key string, duty float32) error {
	if duty < 0 || duty > 100 {
		return fmt.Errorf("invalid duty cycle: %f", duty)
	}

	pwm, ok := exported[key]
	if !ok {
		return fmt.Errorf("PWM '%s' not found", key)
	}

	pwm.duty = uint(float32(pwm.periodNs) * duty * 0.01)

	_, err := fmt.Fprintf(pwm.dutyFile, "%d", pwm.duty)
	return err
}

func Start(key string, duty, freq float32, polarity int) error {
	err := initialize()
	if err != nil {
		return err
	}

	err = bbio.LoadDeviceTree("bone_pwm_" + key)
	if err != nil {
		return err
	}

	//finds and builds the pwmTestPath, as it can be variable...
	pwmTestPath, err := bbio.BuildPath(ocpDir, "pwm_test_"+key)
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

	exported[key] = &export{
		periodFile:   periodFile,
		dutyFile:     dutyFile,
		polarityFile: polarityFile,
	}

	return nil
}

func Disable(key string) {
	bbio.UnloadDeviceTree("bone_pwm_" + key)
	if pwm, ok := exported[key]; ok {
		pwm.periodFile.Close()
		pwm.dutyFile.Close()
		pwm.polarityFile.Close()
		delete(exported, key)
	}
}

func Cleanup() {
	for key := range exported {
		Disable(key)
	}
}
