# notify-brevisone

Notification Plugin for brevis.one SMS Gateways to send SMS messages or ring a contact.

Documentation for your gateway can be found here [docs.brevis.one](https://docs.brevis.one/current/en/Content/Home.htm).

## Installation

Download the current binary from GitHub and install it with Icinga 2. See [releases](https://github.com/NETWAYS/notify-brevisone/releases).

```bash
chmod 755 notify-brevisone
cp notify-brevisone /etc/icinga2/scripts/

/etc/icinga2/scripts/notify-brevisone --help
```

You can use the provided configuration in the `contrib` directory for an example `NotificationCommand`.

```bash
curl -L -o /etc/icinga2/conf.d/brevisone.conf \
https://raw.githubusercontent.com/NETWAYS/notify-brevisone/main/contrib/icinga2/commands.conf
```

You can adjust the command to your needs.

## Usage

```bash
Arguments:
  -g, --gateway string       IP/address of the brevis.one gateway (required)
  -u, --username string      API user name (required)
  -p, --password string      API user password (required)
      --insecure             Skip verification of the TLS certificates (is needed for the default self signed certificate)
      --no-tls                Do NOT use TLS to connect to the gateway (default false)
      --use-legacy-http-api   Use old HTTP API (required on older firmware versions, default false)
  -T, --target string        Target contact, group or phone number (required)
      --target-type string   Target type, one of: number, contact or contactgroup (default "number")
  -R, --ring                 Add ring mode (also ring the target after sending SMS)
      --type string          Icinga $notification.type$ (required)
  -H, --host string          Icinga $host.name$ (required)
  -S, --service string       Icinga $service.name$ (required for service notifications)
  -s, --state string         Icinga $host.state$ or $service.state$ (required)
  -o, --output string        Icinga $host.output or $service.output$ (required)
  -C, --comment string       Icinga $notification.comment$ (optional)
  -a, --author string        Icinga $notification.author$ (optional)
      --date string          Icinga $icinga.long_date_time$ (optional)
  -t, --timeout int          Abort the check after n seconds (default 30)
  -d, --debug                Enable debug mode
  -v, --verbose              Enable verbose mode
  -V, --version              Print version and exit
```

To use this plugin the brevis.one gateway has to be configured properly:

- The REST-API must be activated
- Contacts and Contactgroups can be used for notifications (see the `targetType` argument) if configured
- The credentials must be given to the plugin (`username` and `password`)
- Ensure sending SMS works on the gateway before trying with the plugin
- The device uses a self-signed certificate per default, if you do not intend to replace it,
  adding the `--insecure` flag allows you to ignore this (although this is not recommended)

## Examples

An example for a (custom) service notification with a preconfigured contact:

```bash
notify-brevisOne \
	'--gateway''192.168.0.2' \
	'--username' 'admin' \
	'--password' 'admin' \
	'--target' 'myUser' \
	'--target-type' 'contact' \
	'--comment' 'asfsdf' \
	'--date' '2021-04-14 15:29:06 +0200' \
	'--host' 'myHost' \
	'--service' 'fake' \
	'--type' 'CUSTOM' \
	'--author' 'aRandomMonitor' \
	'--output' 'Hello World' \
	'--state' 'WARNING'
```

This sends the message: `2021-04-14 15:29:06 +0200/CUSTOM: fake @ myHost - WARNING "asfsdf" by aRandomMonitor Hello World`

An example for a (recovery) Host notification with a phone number directly:

```bash
notify-brevisOne \
	'--gateway''192.168.0.2' \
	'--username' 'admin' \
	'--password' 'admin' \
	'--target' '01189998819991197253' \
	'--date' '2021-04-14 15:29:06 +0200' \
	'--host' 'myHost' \
	'--type' 'RECOVERY' \
	'--output' 'It pings again!' \
	'--state' 'OK'
```

This sends the message: `2021-04-14 15:29:06 +0200/RECOVERY: myHost - OK It pings again!`

## License

brevis.one is a trademark of BASIS Europe GmbH

Copyright (C) 2021 [NETWAYS GmbH](mailto:info@netways.de)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
