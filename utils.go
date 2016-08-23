// Package utils implements various global utility functions
package utils

import (
	"bufio"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// Normalize a port string
func NormalizePortString(port string) string {
	return ":" + strings.TrimPrefix(port, ":")
}

// Checks if the port string is valid and returns a normalized version
// if it is
func IsValidPortString(port string) (bool, string) {
	if port[0] == ':' {
		port = strings.TrimPrefix(port, ":")
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return false, ""
	}
	if p > 65535 {
		return false, ""
	}
	return true, ":" + port
}

// Substitures the tilde (~) character for the home directory of the
// current user (but only if its the first character of the string)
func SubstituteHomeDir(path string) string {
	u, err := user.Current()
	if err != nil {
		log.Fatal("Unable to rerieve the current user's information.")
	}
	homeDir := u.HomeDir
	if homeDir == "" {
		log.Fatal("Unable to find the home directory of the current user.")
	}
	if path[:2] == "~/" {
		//path = strings.Replace(path, "~/", "", 1)
		path = filepath.Join(homeDir, path[2:])
	}
	return filepath.Clean(path)
}

// Read a file and return every line as a key/value string array
func ReadKeyValueFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening the file (%s) [%s]", filename, err)
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Read a key=value file and returns a collection of strings indexed by their
// keys
func ParseKeyValueFile(filename string) (ret map[string]string, err error) {
	lines, err := ReadKeyValueFile(filename)
	if err != nil {
		log.Fatalf("Error opening the file (%s) [%s]", filename, err)
		return ret, err
	}
	ret = make(map[string]string)
	for _, line := range lines {
		if strings.Compare(line, "") == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		split := strings.Split(line, "=")
		if len(split) == 2 {
			ret[strings.TrimSpace(split[0])] = strings.TrimSpace(split[1])
		}
	}
	return ret, err
}
