package notes

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

const BUFSIZE int64 = 256

func ReadHeadings(filepath string, numOfHeadings int, level int) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	fi, err := file.Stat()
	if err != nil {
		return "", err
	}
	defer file.Close()

	heading := strings.Repeat("#", level)
	sep := fmt.Sprintf("\n%s ", heading)
	filesize := fi.Size()
	offset := BUFSIZE
	count := 0
	out := ""
	readBuffer := make([]byte, BUFSIZE)
	overflow := ""
	for offset < filesize {
		_, err = file.Seek(-offset, io.SeekEnd)
		if err != nil {
			return "", err
		}
		bytesRead, err := file.Read(readBuffer)
		if err != nil {
			return "", err
		}
		if bytesRead > 0 {
			datas := fmt.Sprint(string(readBuffer), overflow)
			matches := strings.Count(datas, sep)
			count += matches
			if count > numOfHeadings {
				targetString := splitAfterN(datas, sep, count-numOfHeadings+2)
				out = fmt.Sprint(sep, targetString, out)
				return out, nil
			} else {
				before, after, exists := strings.Cut(datas, sep)
				if exists {
					out = fmt.Sprint(sep, after, out)
				}
				overflow = before
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

func splitAfterN(s, sep string, n int) string {
	split := strings.SplitN(s, sep, n)
	return split[len(split)-1]
}

func App(w io.Writer, config Config) {
	content, err := ReadHeadings(config.Filepath, config.NumOfHeadings, config.Level)
	if err != nil {
		log.Fatal(err.Error())
	}
	// _ = content
	Preview(w, content)
}

func Preview(w io.Writer, in string) error {
	out, err := glamour.Render(in, "dark")
	fmt.Fprintln(w, out)
	if err != nil {
		return err
	}
	return nil
}
