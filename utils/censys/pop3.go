package censys

// Represents a dumped Censys POP3 Entry
type POP3Entry struct {
	*BasicEntry
	Banner string `json:"banner"`
}

func (e *POP3Entry) GetBanner() string {
	return e.Banner
}
