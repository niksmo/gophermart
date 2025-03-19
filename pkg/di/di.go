package di

type Row interface {
	Scan(dest ...any) error
}

type Validator interface {
	Validate() (result any, ok bool)
}
