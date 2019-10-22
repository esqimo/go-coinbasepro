package coinbasepro

type Error struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
