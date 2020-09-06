package main

import (
	"errors"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/nsf/termbox-go"
)

// TerminalWidth determines width of current terminal in characters
func TerminalWidth() (int, error) {
	if err := termbox.Init(); err != nil {
		return 0, err
	}
	w, _ := termbox.Size()
	termbox.Close()
	return w, nil
}

// isTerminal checks whether file is a terminal
func isTerminal(file *os.File) bool {
	return isatty.IsTerminal(file.Fd()) || isatty.IsCygwinTerminal(file.Fd())
}

// ParseCookies parses cookies in raw string format into net/http format
func ParseCookies(cookies string) (cookSlice []*http.Cookie, err error) {
	// strip quotes if any
	cookies = strings.Trim(cookies, `"'`)

	// split several cookies into slice
	cookieS := strings.Split(cookies, ";")

	for _, c := range cookieS {
		// strip whitespace
		c = strings.TrimSpace(c)

		// split to name and value
		nameVal := strings.SplitN(c, "=", 2)
		if len(nameVal) != 2 {
			return nil, errors.New("failed to parse cookie")
		}

		cookSlice = append(cookSlice, &http.Cookie{Name: nameVal[0], Value: nameVal[1]})
	}
	return cookSlice, nil
}

// DetectContentType detects HTTP content type based on provided POST data
func DetectContentType(data string) string {
	var contentType string

	if data[0] == '{' || data[0] == '[' {
		contentType = "application/json"
	} else {
		match, _ := regexp.MatchString("([^=]+=[^=]+&?)+", data)
		if match {
			contentType = "application/x-www-form-urlencoded"
		} else {
			contentType = http.DetectContentType([]byte(data))
		}
	}
	return contentType
}