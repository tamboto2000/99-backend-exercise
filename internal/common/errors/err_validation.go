package errors

type Fields map[string][]string

func (f Fields) Add(name, msg string) {
	msgs := f[name]
	msgs = append(msgs, msg)
	f[name] = msgs
}

func (f Fields) NotEmpty() bool {
	return len(f) > 0
}

type ErrValidation struct {
	Err
	fields Fields
}

func NewErrValidation(msg string, fields Fields) ErrValidation {
	return ErrValidation{
		Err: Err{
			msg:  msg,
			code: CodeValidation,
		},
		fields: fields,
	}
}

func (eVld ErrValidation) Fields() Fields {
	return eVld.fields
}

type ErrNotExists struct {
	Err
}

func NewErrNotExists(msg string) ErrNotExists {
	return ErrNotExists{
		Err: Err{
			msg:  msg,
			code: CodeNotExists,
		},
	}
}
