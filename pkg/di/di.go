package di

type Row interface {
	Scan(dest ...any) error
}
