package keys

import "github.com/ygrebnov/keys"

var (
	DBDriver    = keys.New("db.driver")
	JournalMode = keys.New("journal.mode")
	Filename    = keys.New("filename")
)
