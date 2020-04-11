package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	ipExtractionRegex = regexp.MustCompile(`(.+):\d{1,5}$`)
)

func extractIPFromAddress(addr string) (string, error) {
	matches := ipExtractionRegex.FindAllStringSubmatch(addr, -1)
	if len(matches) > 0 && len(matches[0]) >= 1 {
		return strings.Trim(matches[0][1], "[]"), nil
	} else {
		return "", fmt.Errorf("failed to extract IP address from addr %s", addr)
	}
}
