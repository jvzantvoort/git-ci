package changelog

type LogEntry struct {
	SHA1                   string
	SHA1Short              string
	Subject                string
	AuthorName             string
	AuthorEmail            string
	AuthorDate             string
	authorDateTimestamp    string
	CommitterName          string
	CommitterDateTimestamp string
	RawBody                string
	Body                   string
}
