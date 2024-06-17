package gerror

type Error struct {
	Code int
	Msg  string
}

func (err Error) Error() string {
	return err.Msg
}
