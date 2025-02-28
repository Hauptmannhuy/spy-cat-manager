package main

type dbError struct {
	reason string
	code   int
}

type serverError struct {
	reason string
	code   int
}

func (sError serverError) Error() string {
	return sError.reason
}

func (dbErr dbError) Error() string {
	return dbErr.reason
}
