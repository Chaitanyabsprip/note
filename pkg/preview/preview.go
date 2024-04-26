package preview

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

const BUFSIZE int64 = 256

type Preview struct {
	out           io.Writer
	Type          string
	NotesPath     string
	NumOfHeadings int
	Level         int
}

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
	Render(p.out, content)
	return nil
}

func GetHeadings(file *os.File, numOfHeadings int, level int) (string, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	filesize := fileInfo.Size()
	heading := strings.Repeat("#", level)
	sep := fmt.Sprintf("\n%s ", heading)
	offset := BUFSIZE
	count := 0
	out := ""
	readBuffer := make([]byte, BUFSIZE)
	overflow := ""
	for offset < filesize {
		var bytesRead int
		bytesRead, err = readBefore(file, offset, readBuffer)
		if err != nil {
			return "", err
		}
		if bytesRead > 0 {
			count, out, overflow = parsePartialHeadings(readBuffer, count, numOfHeadings, sep, out, overflow)
			if count == numOfHeadings {
				return out, nil
			}
		}
		offset += BUFSIZE
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}
	bytesRead, err := file.Read(readBuffer)
	if err != nil {
		return "", err
	}
	if bytesRead > 0 {
		readBuffer = readBuffer[:filesize-offset+BUFSIZE]
		out = fmt.Sprint(string(readBuffer), overflow, out)
	}
	return out, nil
}

func readBefore(file io.ReadSeeker, offset int64, readBuffer []byte) (int, error) {
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
	readBuffer []byte,
	count, numOfHeadings int,
	sep,
	out,
	overflow string,
) (int, string, string) {
	datas := fmt.Sprint(string(readBuffer), overflow)
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

func Render(w io.Writer, in string) error {
	out, err := glamour.Render(in, "dark")
	fmt.Fprintln(w, out)
	if err != nil {
		return err
	}
	return nil
}
