package issue

import "time"

type Issue struct {
	ID          string
	Title       string
	Description string
	Repository  Repository
	Labels      []string
	Comments    []Comment
	Status      Status
}

type Status int

const (
	Open Status = iota + 1
	Closed
	InProgress
)

type Comment struct {
	CreatedAt time.Time
	ID        string
	Body      string
	Issue     Issue
}

type Repository struct {
	ID   string
	Name string
	Path string
}

func New(id, title, desc string, labels []string) *Issue {
	issue := new(Issue)
	return issue
}

func Update(id, title, desc string, labels []string, status Status) *Issue {
	issue := new(Issue) // get old issue from db/file
	// update issue
	return issue
}

func Find(query string) *Issue {
	issue := new(Issue) // get old issue from db/file
	return issue
}

func (i *Issue) NewComment() {
}
