package git

var (
	GitFormatKeys = map[string]string{
		"sha1":                     "%H",
		"sha1_short":               "%h",
		"subject":                  "%s",
		"author_name":              "%an",
		"author_email":             "%ae",
		"author_date":              "%ad",
		"author_date_timestamp":    "%at",
		"committer_name":           "%cn",
		"committer_date_timestamp": "%ct",
		"raw_body":                 "%B",
		"body":                     "%b",
	}
)

func (g *GitCmd) Initialize() {
	for keyname, value := range GitFormatKeys {
		g.FormatKeyNames = append(g.FormatKeyNames, keyname)
		g.FormatKeyValues = append(g.FormatKeyValues, value)
	}
}


