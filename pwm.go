package bbio

import (
	"fmt"
	"os"
)

type PWM struct {
	key          string
	dutyCycle    float32
	frequency    float32
	polarity     int
	periodFile   *os.File
	dutyFile     *os.File
	polarityFile *os.File
}

var (
	pwmInitialized bool
)

func NewPWM(key string, duty, frequency float32, polarity int) (*PWM, error) {
	if !pwmInitialized {
		err := LoadDeviceTree("am33xx_pwm")
		if err != nil {
			return nil, err
		}
		pwmInitialized = true
	}

	err := LoadDeviceTree("bone_pwm_" + key)
	if err != nil {
		return nil, err
	}

	ocpDir, err := BuildPath("/sys/devices", "ocp")
	if err != nil {
		return nil, err
	}

	//finds and builds the pwmTestPath, as it can be variable...
	pwmTestPath, err := BuildPath(ocpDir, "pwm_test_"+key)
	if err != nil {
		return nil, err
	}

	//create the path for the period and duty
	periodPath := pwmTestPath + "/period"
	dutyPath := pwmTestPath + "/duty"
	polarityPath := pwmTestPath + "/polarity"

	periodFile, err := os.OpenFile(periodPath, os.O_RDWR, 0660)
	if err != nil {
		return nil, err
	}
	dutyFile, err := os.OpenFile(dutyPath, os.O_RDWR, 0660)
	if err != nil {
		periodFile.Close()
		return nil, err
	}
	polarityFile, err := os.OpenFile(polarityPath, os.O_RDWR, 0660)
	if err != nil {
		periodFile.Close()
		dutyFile.Close()
		return nil, err
	}

	pwm := &PWM{
		key:          key,
		periodFile:   periodFile,
		dutyFile:     dutyFile,
		polarityFile: polarityFile,
	}

	err = pwm.SetFrequency(frequency)
	if err != nil {
		pwm.Close()
		return nil, err
	}
	err = pwm.SetPolarity(polarity)
	if err != nil {
		pwm.Close()
		return nil, err
	}
	err = pwm.SetDutyCycle(duty)
	if err != nil {
		pwm.Close()
		return nil, err
	}

	return pwm, nil
}

func (pwm *PWM) Key() string {
	return pwm.Key()
}

func (pwm *PWM) Frequency() float32 {
	return pwm.frequency
}

func (pwm *PWM) SetFrequency(frequency float32) error {
	if frequency <= 0 {
		return fmt.Errorf("invalid requency: %f", frequency)
	}
	if frequency == pwm.frequency {
		return nil
	}

	periodNs := uint(1e9 / frequency)
	_, err := fmt.Fprintf(pwm.periodFile, "%d", periodNs)
	if err != nil {
		return err
	}

	pwm.frequency = frequency
	return nil
}

func (pwm *PWM) Polarity() int {
	return pwm.polarity
}

func (pwm *PWM) SetPolarity(polarity int) error {
	if polarity == pwm.polarity {
		return nil
	}

	_, err := fmt.Fprintf(pwm.polarityFile, "%d", polarity)
	if err != nil {
		return err
	}

	pwm.polarity = polarity
	return nil
}

func (pwm *PWM) DutyCycle() float32 {
	return pwm.dutyCycle
}

func (pwm *PWM) SetDutyCycle(dutyCycle float32) error {
	if dutyCycle < 0 || dutyCycle > 100 {
		return fmt.Errorf("invalid duty cycle: %f", dutyCycle)
	}
	if dutyCycle == pwm.dutyCycle {
		return nil
	}

	periodNs := uint(1e9 / pwm.frequency)
	duty := uint(float32(periodNs) * dutyCycle * 0.01)
	_, err := fmt.Fprintf(pwm.dutyFile, "%d", duty)
	if err != nil {
		return err
	}

	pwm.dutyCycle = dutyCycle
	return nil
}

func (pwm *PWM) Close() {
	UnloadDeviceTree("bone_pwm_" + pwm.key)
	pwm.periodFile.Close()
	pwm.dutyFile.Close()
	pwm.polarityFile.Close()
}

func CleanupPWM() error {
	pwmInitialized = false
	return UnloadDeviceTree("am33xx_pwm")
}
