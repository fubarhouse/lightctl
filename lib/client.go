package lib

import "net"

func NewClient() *client {
	return &client{}
}

func (c *client) SetRequest(request APIRequest) {
	c.Request = &request
}

func (c *client) SetResponse(response APIResponse) {
	c.Response = &response
}

func (c *client) SetParameters(properties ValueProperties) {
	c.Properties.RequestModifier = properties
}

func (c *client) SetIPs(value []net.IP) {
	c.Properties.IPs = value
}

func (c *client) SetPort(value string) {
	c.Properties.Port = value
}
