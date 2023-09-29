package department

// Teacher represents a teacher.
type Teacher struct {
	Username string
}

// GetWebsite returns the website of the teacher.
func (t Teacher) GetWebsite() string {
	return "https://www.unibo.it/sitoweb/" + t.Username
}
