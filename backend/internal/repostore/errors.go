package repostore

import "fmt"

var (
	ErrRecordNotFound      = fmt.Errorf("no such record")
	ErrRollbackTransaction = fmt.Errorf("unable to rollback transaction")
	ErrForbiddenParameters = fmt.Errorf("recieved forbidden parameters")
)
