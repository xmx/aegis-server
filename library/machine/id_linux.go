package machine

// machineID returns the uuid specified at `/var/lib/dbus/machine-id` or `/etc/machine-id`.
// If there is an error reading the files an empty string is returned.
// See https://unix.stackexchange.com/questions/144812/generate-consistent-machine-unique-id
func machineID() (string, error) {
	mid, err := readFile("/etc/machine-id")
	if err != nil {
		mid, err = readFile("/var/lib/dbus/machine-id")
	}

	return mid, err
}
