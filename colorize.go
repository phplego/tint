package tint

import (
	"io"
	"regexp"
	"strings"
)

var (
	colors = map[string]string{
		"k": "30", "K": "90",
		"r": "31", "R": "91",
		"g": "32", "G": "92",
		"y": "33", "Y": "93",
		"b": "34", "B": "94",
		"m": "35", "M": "95",
		"c": "36", "C": "96",
		"w": "37", "W": "97",
	}

	bgColors = map[string]string{
		"k": "40", "K": "100",
		"r": "41", "R": "101",
		"g": "42", "G": "102",
		"y": "43", "Y": "103",
		"b": "44", "B": "104",
		"m": "45", "M": "105",
		"c": "46", "C": "106",
		"w": "47", "W": "107",
	}

	re      = regexp.MustCompile(`(?s)@[kKrRgGyYbBmMcCwW*]!?[kKrRgGyYbBmMcCwW]?{.*?}`)
	reStrip = regexp.MustCompile(`\033\[\d+(;\d+)*m`)
)

type ColorizeWriter struct {
	out io.Writer
}

func (w ColorizeWriter) Write(p []byte) (n int, err error) {
	return w.out.Write([]byte(Colorize(string(p))))
}

func strip(s string) string {
	return reStrip.ReplaceAllString(s, "")
}

func applyColor(text string, format string) string {
	result := strings.Builder{}
	result.WriteString("\033[")
	if strings.Contains(format, "!") {
		result.WriteString("01;")
		format = strings.ReplaceAll(format, "!", "")
	}
	result.WriteString(colors[format[0:1]])
	if len(format) == 2 {
		result.WriteString(";")
		result.WriteString(bgColors[format[1:2]])
	}
	result.WriteString("m")
	result.WriteString(text)
	result.WriteString("\033[00m")
	return result.String()
}

func applyRainbow(text string) string {
	rainbow := "RYGCBM"
	result := strings.Builder{}
	var i int
	for _, c := range text {
		if c == ' ' {
			result.WriteRune(c)
			continue
		}
		pos := i % len(rainbow)
		color := colors[rainbow[pos:pos+1]]
		result.WriteString("\033[01;" + color + "m" + string(c) + "\033[00m")
		i++
	}
	return result.String()
}

func Colorize(s string) string {
	return re.ReplaceAllStringFunc(s, func(m string) string {
		format := m[1:strings.Index(m, "{")]
		content := m[strings.Index(m, "{")+1 : len(m)-1]

		if strings.Contains(format, "*") {
			return applyRainbow(content)
		}
		return applyColor(content, format)
	})
}
