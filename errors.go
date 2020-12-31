package main

import (
	"fmt"

	"golang.org/x/xerrors"
)

// ErrInvalidMessage is Messageが処理できない時に返す
var ErrInvalidMessage = &Error{
	Code:    "InvalidMessage",
	Message: "InvalidMessage",
	KV:      map[string]interface{}{},
}

// Error is Error情報を保持する struct
type Error struct {
	Code    string
	Message string
	KV      map[string]interface{}
	err     error
}

// Error is error interface func
func (e *Error) Error() string {
	if e.KV == nil || len(e.KV) < 1 {
		return fmt.Sprintf("%s: %s: %s", e.Code, e.Message, e.err)
	}
	return fmt.Sprintf("%s: %s: attribute:%+v :%s", e.Code, e.Message, e.KV, e.err)
}

// Is is err equal check
func (e *Error) Is(target error) bool {
	var appErr *Error
	if !xerrors.As(target, &appErr) {
		return false
	}
	return e.Code == appErr.Code
}

// Unwrap is return unwrap error
func (e *Error) Unwrap() error {
	return e.err
}

// newErrInvalidMessage is return ErrInvalidMessage
func newErrInvalidMessage(message string, kv map[string]interface{}, err error) *Error {
	return &Error{
		Code:    ErrInvalidMessage.Code,
		Message: message,
		KV:      kv,
		err:     err,
	}
}
