// Package preview provides preview  
package preview

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

// BUFSIZE  
const BUFSIZE int64 = 256

// Preview struct  
type Preview struct {
	out           io.Writer
	Type          string
	NotesPath     string
	NumOfHeadings int
	Level         int
}

// New function  
func New(
	w io.Writer,
	_type string,
	notesPath string,
	numOfHeadings int,
	level int,
) *Preview {
	p := new(Preview)
	p.Type = _type
	p.NotesPath = notesPath
	p.NumOfHeadings = numOfHeadings
	p.Level = level
	p.out = w
	return p
}

// Peek method  
func (p *Preview) Peek() error {
	file, err := os.Open(p.NotesPath)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := GetHeadings(file, p.NumOfHeadings, p.Level)
	if err != nil {
		return err
	}
	err = Render(p.out, content)
	if err != nil {
		return err
	}
	return nil
}

// GetHeadings function  
func GetHeadings(file *os.File, numOfHeadings int, level int) (string, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	filesize := fileInfo.Size()
	heading := strings.Repeat("#", level)
	sep := fmt.Sprintf("\n%s ", heading)
	var prevOffset int64
	offset := BUFSIZE
	count := 0
	out := ""
	readBuffer := make([]byte, BUFSIZE)
	overflow := ""
	for offset < filesize+BUFSIZE {
		var bytesRead int
		bytesRead, err = readBefore(file, min(offset, filesize), readBuffer)
		if err != nil {
			return "", err
		}
		if bytesRead > 0 {
			readString := string(readBuffer[:min(filesize, offset)-prevOffset])
			count, out, overflow = parsePartialHeadings(
				readString,
				sep,
				out,
				overflow,
				count,
				numOfHeadings,
			)
			if count == numOfHeadings {
				return out, nil
			}
		}
		prevOffset = offset
		offset += BUFSIZE
	}
	return out, nil
}

func readBefore(
	file io.ReadSeeker,
	offset int64,
	readBuffer []byte,
) (int, error) {
	_, err := file.Seek(-offset, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	var bytesRead int
	bytesRead, err = file.Read(readBuffer)
	if err != nil {
		return 0, err
	}
	return bytesRead, nil
}

func parsePartialHeadings(
	readString,
	sep,
	out,
	overflow string,
	count, numOfHeadings int,
) (int, string, string) {
	datas := fmt.Sprint(readString, overflow)
	matches := strings.Count(datas, sep)
	count += matches
	if count > numOfHeadings {
		targetString := splitAfterN(datas, sep, count-numOfHeadings+2)
		out = fmt.Sprint(sep, targetString, out)
		return numOfHeadings, out, overflow
	}
	before, after, exists := strings.Cut(datas, sep)
	if exists {
		out = fmt.Sprint(sep, after, out)
	}
	overflow = before
	return count, out, overflow
}

func splitAfterN(s, sep string, n int) string {
	split := strings.SplitN(s, sep, n)
	return split[len(split)-1]
}

// Render function  
func Render(w io.Writer, in string) error {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(glamour.DarkStyle),
		glamour.WithWordWrap(120),
		glamour.WithPreservedNewLines(),
	)
	if err != nil {
		return err
	}
	out, err := renderer.Render(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, out)
	return nil
}
