# Lightctl

**A command-line control program for Elgato Key Lights running on your local network**

## Special thanks

It wasn't as difficult as I thought initially, however I want to give special thanks to the following resources:

- [Controlling Elgato Key Light under Ubuntu with Ruby](https://mensfeld.pl/2021/12/controlling-elgato-key-light-under-ubuntu-with-ruby/)
- [Elgato Key Light API](https://github.com/adamesch/elgato-key-light-api)

## Installation

Installation docs are on the way, but this will definitely make its way to the AUR.

## Usage

**Flags**

| Flag                                  | Description                                      |
|---------------------------------------|--------------------------------------------------|
| `--ip 192.168.0.60 --ip 192.168.0.61` | IP addresses of your lights                      |
| `--port 666`                          | Port of the exposed IP addresses                 |
| `--value +20`                         | Value in which to set, increment or decrement by |

**Note**: Decreasing a value is currently not working.

**Commands**

| Command                          | Description                                                                                                                            |
|----------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|
| `lightctl on`                    | Turns each light on.                                                                                                                   |
| `lightctl off`                   | Turns each light off.                                                                                                                  |
| `lightctl toggle`                | Toggles each light to an on or off state.                                                                                              |
| `lightctl info`                  | Show some information from the `accessory-info` uri on the API.                                                                        |
| `lightctl state`                 | Show some information from the `light` usi on the endpoint                                                                             |
| `lightctl brightness --value x`  | Set the value of brightness to a specific value such as `20`, `-10` or `+10`<br />Note: Negative value input does not work right now.  |
| `lightctl temperature --value x` | Set the value of temperature to a specific value such as `20`, `-10` or `+10`<br />Note: Negative value input does not work right now. |


## Known issues

- Setting a value intended to decrease the current value results in a failure and the device turns off.

## License

MIT - no obligations or warranties are provided with this application.