package term

import (
	"os"
	"strings"

	"golang.org/x/term"
)

func IsTerminal(f ...*os.File) bool {
	for _, file := range f {
		if !term.IsTerminal(int(file.Fd())) {
			return false
		}
	}
	if len(f) == 0 {
		return term.IsTerminal(int(stdin))
	}
	return true
}

func IsPipe(f ...*os.File) bool {
	for _, file := range f {
		if !isPipe(file) {
			return false
		}
	}
	if len(f) == 0 {
		return isPipe(os.Stdin)
	}
	return true
}

func IsFile(f ...*os.File) bool {
	for _, file := range f {
		if !isFile(file) {
			return false
		}
	}
	if len(f) == 0 {
		return isFile(os.Stdin)
	}
	return true
}

func IsDumb() bool {
	return os.Getenv("TERM") == "dumb"
}

func IsWSL() bool {
	if os.Getenv("WSL_INTEROP") != "" ||
		os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}

	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}

	v := strings.ToLower(string(data))
	return strings.Contains(v, "microsoft")
}

func IsSSH() bool {
	return os.Getenv("SSH_CONNECTION") != "" ||
		os.Getenv("SSH_CLIENT") != "" ||
		os.Getenv("SSH_TTY") != ""
}

func IsDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	data, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}

	s := string(data)

	return strings.Contains(s, "docker") ||
		strings.Contains(s, "containerd") ||
		strings.Contains(s, "kubepods")
}

func IsRaw() bool {
	return rawModeState != nil
}

func SupportsColor() bool {
	return colorLevel != ColorModeNone
}

func SupportsColorMode(mode ColorMode) bool {
	return colorLevel >= mode
}

func isPipe(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeNamedPipe != 0
}

func isFile(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}
