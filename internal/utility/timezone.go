package utility

// Source - https://stackoverflow.com/a/78745212
// Posted by clubcleaver
// Retrieved 2026-05-05, License - CC BY-SA 4.0

import (
	"os"
	"runtime"
	"strings"
	"time"
)

var zoneDirs = map[string]string{
	"android":   "/system/usr/share/zoneinfo/",
	"darwin":    "/usr/share/zoneinfo/",
	"dragonfly": "/usr/share/zoneinfo/",
	"freebsd":   "/usr/share/zoneinfo/",
	"linux":     "/usr/share/zoneinfo/",
	"netbsd":    "/usr/share/zoneinfo/",
	"openbsd":   "/usr/share/zoneinfo/",
	"solaris":   "/usr/share/lib/zoneinfo/",
}

var timeZones []string

func GetAllTimezones(useCache bool) ([]string, error) {
	if len(timeZones) > 0 && useCache {
		return timeZones, nil
	}

	// Reads the Directory corresponding to the OS
	dirFile, _ := os.ReadDir(zoneDirs[runtime.GOOS])
	for _, i := range dirFile {
		// Checks if starts with Capital Letter
		if i.Name() == (strings.ToUpper(i.Name()[:1]) + i.Name()[1:]) {
			if i.IsDir() {
				// Recursive read if directory
				subFiles, err := os.ReadDir(zoneDirs[runtime.GOOS] + i.Name())
				if err != nil {
					return nil, err
				}
				for _, s := range subFiles {
					// Appends the path to timeZones var
					timeZones = append(timeZones, i.Name()+"/"+s.Name())
				}
			}
			// Appends the path to timeZones var
			timeZones = append(timeZones, i.Name())
		}
	}
	// Loop over timezones and Check Validity, Delete entry if invalid.
	// Range function doesnt work with changing length.
	for i := 0; i < len(timeZones); i++ {
		_, err := time.LoadLocation(timeZones[i])
		if err != nil {
			// newSlice = timeZones[:n]  timeZones[n+1:]
			timeZones = append(timeZones[:i], timeZones[i+1:]...)
			continue
		}
	}

	return timeZones, nil
}
