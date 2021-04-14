# notify-brevisOne
Notification Plugin for Brevis One SMS Gateways

# Preparation

## Building
After cloning the repository, execute the following in the directory
```
go build
```
## Installation

 - Copy `notify-brevisOne` to the target directory (`/etc/icinga2/scripts` for example) and fix the permissions as necessary (chmod 755 ...).
 - Copy the `NotificationCommand` and the `Notification` templates configuration from the `icinga2_conf` directory into your icinga2 configuration. This is not strictly necessary, you can create them yourself, if you know how to.

# Usage
A call of
```
./notify-brevisOne -h
```
will show the available options. To use this plugin the BrevisOne gateway has to be configured beforehand.

 - The REST-API must be activated
 - Contacts and Contactgroups can be used for notifications (see the `targetType` argument) if configured
 - The credentials must be given to the plugin (`username` and `password`)
 - SIM card must be available and the functionality tested
 - The device uses a self-signed certificate per default, if you do not intend to fix this, the `--skipTlsVerify` flag allows you to ignore this (although this is not recommended)

## Examples
An example for a (custom) service notification with a preconfigured contact:
```
./notify-brevisOne \
	'--gateway''192.168.0.2' \
	'--username' 'admin' \
	'--password' 'admin' \
	'--target' 'myUser' \
	'--targetType' 'contact' \
	'--comment' 'asfsdf' \
	'--date' '2021-04-14 15:29:06 +0200' \
	'--host' 'myHost' \
	'--service' 'fake' \
	'--type' 'CUSTOM' \
	'-a' 'aRandomMonitor' \
	'-o' 'Hello World' \
	'--checkresult' 'WARNING'
```

An example for a (recovery) Host notification with a phone number directly
```
./notify-brevisOne \
	'--gateway''192.168.0.2' \
	'--username' 'admin' \
	'--password' 'admin' \
	'--target' '01189998819991197253' \
	'--date' '2021-04-14 15:29:06 +0200' \
	'--host' 'myHost' \
	'--type' 'RECOVERY' \
	'-a' 'aRandomMonitor' \
	'-o' 'It pings again!' \
	'--checkresult' 'OK'
```
