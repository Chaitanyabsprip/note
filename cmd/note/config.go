package main

type Config struct {
	Content         string
	Type            string
	Notespath       string
	defaultFilename string
	Level           int
	NumOfHeadings   int
	EditFile        bool
	Global          bool
	Quiet           bool
}
