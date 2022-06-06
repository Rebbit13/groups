package service

type BasicError struct {
	message string
}

func (be *BasicError) Error() string {
	return be.message
}

type GroupToAttachNotExistError struct {
	BasicError
}

func NewGroupToAttachNotExistError(message string) *GroupToAttachNotExistError {
	return &GroupToAttachNotExistError{BasicError: BasicError{message: message}}
}

type RecursiveGroupDependenciesError struct {
	BasicError
}

func NewRecursiveGroupDependenciesError(message string) *RecursiveGroupDependenciesError {
	return &RecursiveGroupDependenciesError{BasicError: BasicError{message: message}}
}

type RecordNotFound struct {
	BasicError
}

func NewRecordNotFound(message string) *RecordNotFound {
	return &RecordNotFound{BasicError: BasicError{message: message}}
}
