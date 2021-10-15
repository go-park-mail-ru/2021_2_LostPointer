package models

type CustomError struct {
	ErrorType 	  int
	OriginalError error
	Message 	  string
}