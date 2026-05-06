package entity

import "errors"

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTaskID     = errors.New("invalid task id")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrInvalidTaskStatus = errors.New("invalid task status")
	ErrMissingTaskID     = errors.New("task id is required")
	ErrMissingTaskTitle  = errors.New("task title is required")
)
