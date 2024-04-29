package note

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chaitanyabsprip/note/pkg/preview"
)

func sentenceCase(input string) string {
	var sb strings.Builder
	sentences := strings.Split(input, ". ")
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) == 0 {
			continue
		}
		sentence = strings.ToLower(sentence)
		sb.WriteString(strings.ToUpper(string(sentence[0])))
		sb.WriteString(sentence[1:])
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

const wrapWidth = 80

func wordWrap(text string, lineWidth int) string {
	lines := strings.Split(text, "\n")
	wrapped := ""
	for _, line := range lines {
		words := strings.Fields(strings.TrimSpace(line))
		if len(words) == 0 {
			wrapped += line + "\n"
			continue
		}
		currLine := words[0]
		for _, word := range words[1:] {
			if len(currLine)+len(word) <= lineWidth-3 {
				currLine += " " + word
			} else {
				wrapped += currLine + "\n"
				currLine = word
			}
		}
		if currLine != "" {
			wrapped += currLine + "\n"
		}
	}
	return wrapped
}

func addHeading(body string, file *os.File) (string, error) {
	heading, err := newHeading(file)
	if err != nil {
		return "", err
	}
	note := body
	if heading != "" {
		note = fmt.Sprintf("\n%s\n\n%s", heading, note)
	}
	return note, nil
}

func newHeading(file *os.File) (string, error) {
	content, err := preview.GetHeadings(file, 1, 2)
	if err != nil {
		return "", err
	}
	lines := strings.Split(content, "\n")
	lHeading := lastHeading(lines)
	prevTime := strings.TrimPrefix(lHeading, "## ")
	currTime := time.Now().Format("Mon, 02 Jan 2006")
	if currTime != prevTime || lHeading == "" {
		return fmt.Sprint("## ", currTime), nil
	}
	return "", nil
}

func lastHeading(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "##") {
			return line
		}
	}
	return ""
}
