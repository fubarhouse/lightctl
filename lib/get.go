package lib

// GetInfo will show the response from an 'accessory-info' API call.
func (c *client) GetInfo(ip string) (APIResponse, error) {
	return c.DispatchResponse(ip, APIRequest{}, "accessory-info", "GET")
}

// GetStatus will return the state property of a light on a given IP.
func (c *client) GetStatus(ip string) int {
	*c.Response, _ = c.State(ip)
	return c.Response.Lights[0].On
}

// GetTemperature will return the temperature property of a light on a given IP.
func (c *client) GetTemperature(ip string) int {
	*c.Response, _ = c.State(ip)
	return c.Response.Lights[0].Temperature
}

// GetBrightness will return the brightness property of a light on a given IP.
func (c *client) GetBrightness(ip string) int {
	*c.Response, _ = c.State(ip)
	return c.Response.Lights[0].Brightness
}
