package metricscollector

type StorageError struct {
	Err error
}

func (e *StorageError) Error() string {
	return e.Err.Error()
}

func NewStorageError(err error) error {
	return &StorageError{
		Err: err,
	}
}

type StorageRetryableError struct {
	Err error
}

func NewStorageRetryableError(err error) error {
	return &StorageRetryableError{
		Err: err,
	}
}

func (e *StorageRetryableError) Error() string {
	return e.Err.Error()
}
