package lib

// SetStatusToggle will toggle a light on a given IP to have the opposite value as currently set.
func (c *client) SetStatusToggle(ip string) (APIResponse, error) {
	currentState := c.GetStatus(ip)
	Request := APIRequest{}

	if currentState == 0 {
		Request = APIRequest{
			NumberOfLights: 1,
			Lights: []Light{{
				On: 1,
			}},
		}
	}
	if currentState == 1 {
		Request = APIRequest{
			NumberOfLights: 1,
			Lights: []Light{{
				On: 0,
			}},
		}
	}

	return c.DispatchResponse(ip, Request, "lights", "PUT")
}

// SetStatusOff will set a light on a given IP to have an off state.
func (c *client) SetStatusOff(ip string) (APIResponse, error) {
	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On: 0,
		}},
	}
	return c.DispatchResponse(ip, Request, "lights", "PUT")
}

// SetStatusOn will set a light on a given IP to have an on state.
func (c *client) SetStatusOn(ip string) (APIResponse, error) {
	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On: 1,
		}},
	}
	return c.DispatchResponse(ip, Request, "lights", "PUT")
}

// SetTemperature will set a light on a given IP to a specific temperature value.
func (c *client) SetTemperature(ip string) (APIResponse, error) {

	if !c.Properties.RequestModifier.IsNeutral {
		t := c.GetTemperature(ip)
		if c.Properties.RequestModifier.IsPositive {
			c.Properties.RequestModifier.Value = add(c.Properties.RequestModifier.Value, t)
		} else if c.Properties.RequestModifier.IsNegative {
			c.Properties.RequestModifier.Value = subtract(c.Properties.RequestModifier.Value, t)
		}
	}

	if c.Properties.RequestModifier.Value < SettingMinimumTemperature || c.Properties.RequestModifier.Value > SettingMaximumTemperature {
		return APIResponse{}, c.ErrTemperatureOuterBounds(c.Properties.RequestModifier.Value)
	}

	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On:          1,
			Temperature: c.Properties.RequestModifier.Value,
		}},
	}

	return c.DispatchResponse(ip, Request, "lights", "PUT")
}

// SetBrightness will set a light on a given IP to a specific brightness value.
func (c *client) SetBrightness(ip string) (APIResponse, error) {

	if !c.Properties.RequestModifier.IsNeutral {
		b := c.GetBrightness(ip)
		if c.Properties.RequestModifier.IsPositive {
			c.Properties.RequestModifier.Value = add(c.Properties.RequestModifier.Value, b)
		} else if c.Properties.RequestModifier.IsNegative {
			c.Properties.RequestModifier.Value = subtract(c.Properties.RequestModifier.Value, b)
		}
	}

	if c.Properties.RequestModifier.Value < SettingMinimumBrightness || c.Properties.RequestModifier.Value > SettingMaximumBrightness {
		return APIResponse{}, c.ErrBrightnessOuterBounds(c.Properties.RequestModifier.Value)
	}

	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On:         1,
			Brightness: c.Properties.RequestModifier.Value,
		}},
	}

	return c.DispatchResponse(ip, Request, "lights", "PUT")
}
