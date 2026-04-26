package errors

type OptimisticLockError struct{}

func NewOptimisticLockError() *OptimisticLockError {
	return &OptimisticLockError{}
}

func (e *OptimisticLockError) Error() string {
	return "optimistic lock failed, retry required"
}

func (e *OptimisticLockError) Is(target error) bool {
	_, ok := target.(*OptimisticLockError)
	return ok
}
