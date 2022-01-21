package lib

import (
	"errors"
	"fmt"
)

// ErrBrightnessOuterBounds will return a custom error message for brightness limitations based on input.
func (c *client) ErrBrightnessOuterBounds(value int) error {
	return errors.New(fmt.Sprintf("Needs to be between %v and %v, was: %v", SettingMinimumBrightness, SettingMaximumBrightness, value))
}

// ErrTemperatureOuterBounds will return a custom error message for temperature limitations based on input.
func (c *client) ErrTemperatureOuterBounds(value int) error {
	return errors.New(fmt.Sprintf("Needs to be between %v and %v, was: %v", SettingMinimumTemperature, SettingMaximumTemperature, value))
}
