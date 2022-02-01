package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/fubarhouse/lightctl/lib"
	"os"
)

var (
	// Thanks of the extra special kind to the following resources for making this possible.
	Thanks = []string{
		"https://mensfeld.pl/2021/12/controlling-elgato-key-light-under-ubuntu-with-ruby/",
		"https://github.com/adamesch/elgato-key-light-api",
	}

	// CLI configuration via Kingpin.
	app = kingpin.New("chat", "A command-line chat application.")

	on          = app.Command("on", "Turn your light(s) on")
	off         = app.Command("off", "Turn your light(s) off")
	toggle      = app.Command("toggle", "Toggle your light(s) on or off")
	info        = app.Command("info", "Display your light(s) information")
	status      = app.Command("state", "Display your light(s) state information")
	brightness  = app.Command("brightness", fmt.Sprintf("Set your light(s) brightness to a specific value or decrement/increment (between %v and %v)", lib.SettingMinimumBrightness, lib.SettingMaximumBrightness))
	temperature = app.Command("temperature", fmt.Sprintf("Set your light(s) temperature to a specific value or decrement/increment (between %v and %v)", lib.SettingMinimumTemperature, lib.SettingMaximumTemperature))

	destinationValue = app.Flag("value", "value to change to. ie. 40 || 600 || _40 || +40").String()
	destinationIPs   = app.Flag("ip", "IP addresses, multiple supported").Default("10.0.0.90", "10.0.0.91").IPList()
	destinationPort  = app.Flag("port", "Port to use, single value support only").Default("9123").String()
)

func main() {

	app.HelpFlag.Short('h')
	client := lib.NewClient()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case on.FullCommand():
		client.SetIPs(*destinationIPs)
		client.SetPort(*destinationPort)
		for _, v := range *destinationIPs {
			_, err := client.SetStatusOn(v.String())
			if err != nil {
				fmt.Println(err)
			}
			client.PrintResponse()
		}

	case off.FullCommand():
		client.SetIPs(*destinationIPs)
		client.SetPort(*destinationPort)
		for _, v := range *destinationIPs {
			_, err := client.SetStatusOff(v.String())
			if err != nil {
				fmt.Println(err)
			}
			client.PrintResponse()
		}

	case toggle.FullCommand():
		client.SetIPs(*destinationIPs)
		client.SetPort(*destinationPort)
		for _, v := range *destinationIPs {
			_, err := client.SetStatusToggle(v.String())
			if err != nil {
				fmt.Println(err)
			}
			client.PrintResponse()
		}

	case info.FullCommand():
		client.SetIPs(*destinationIPs)
		client.SetPort(*destinationPort)
		for _, v := range *destinationIPs {
			_, err := client.GetInfo(v.String())
			if err != nil {
				fmt.Println(err)
			}
			client.PrintResponse()
		}

	case status.FullCommand():
		client.SetIPs(*destinationIPs)
		client.SetPort(*destinationPort)
		for _, v := range *destinationIPs {
			_, err := client.State(v.String())
			if err != nil {
				fmt.Println(err)
			}
			client.PrintResponse()
		}

	case brightness.FullCommand():
		client.SetIPs(*destinationIPs)
		client.SetPort(*destinationPort)
		for _, v := range *destinationIPs {
			err := client.ParseValue(*destinationValue)
			if err != nil {
				fmt.Println(err)
			}
			_, err = client.SetBrightness(v.String())
			if err != nil {
				fmt.Println(err)
			}
			client.PrintResponse()
		}

	case temperature.FullCommand():
		client.SetIPs(*destinationIPs)
		client.SetPort(*destinationPort)
		for _, v := range *destinationIPs {
			err := client.ParseValue(*destinationValue)
			if err != nil {
				fmt.Println(err)
			}
			_, err = client.SetTemperature(v.String())
			if err != nil {
				fmt.Println(err)
			}
			client.PrintResponse()
		}
	}
}
