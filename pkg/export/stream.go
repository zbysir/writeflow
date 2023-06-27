package export

type Stream interface {
	NewReader() Reader
}

type Reader interface {
	Read() (string, error)
	ReadAll() ([]string, error)
}
