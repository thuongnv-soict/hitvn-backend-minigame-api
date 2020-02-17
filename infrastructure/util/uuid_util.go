package util

import uuid "github.com/nu7hatch/gouuid"


/**
 * Returns a new Uuid
 */
func NewUuid() string {
	var uuid_ *uuid.UUID
	uuid_,_ = uuid.NewV4()

	return (*uuid_).String()
}