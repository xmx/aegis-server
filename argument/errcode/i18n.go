package errcode

func NewI18nError(code int, key string) I18nError {
	return I18nError{
		Code: code,
		Key:  key,
	}
}

type I18nError struct {
	Code int
	Key  string
	Args []any
}

func (e I18nError) Error() string {
	return e.Key
}

func (e I18nError) Fmt(args ...any) I18nError {
	e.Args = args
	return e
}
