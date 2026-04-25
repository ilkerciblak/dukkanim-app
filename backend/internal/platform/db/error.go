package db

import (
	"errors"

	"github.com/lib/pq"
)

const (
 PgUniqueViolation        = "23505"
 PgForeignKeyViolation    = "23503"
 PgNotNullViolation       = "23502"
 PgCheckViolation         = "23514"
 PgRestrictViolation      = "23001"
 PgSerializationFailure   = "40001"
 PgDeadlockDetected       = "40P01"
 PgQueryCanceled          = "57014"
 PgConnectionDoesNotExist = "08003"
 PgConnectionFailure      = "08006"
 PgUndefinedTable         = "42P01"
 PgUndefinedColumn        = "42703"
 PgTooManyConnections     = "53300"
 PgStringTruncation       = "22001"
 PgNumericOutOfRange      = "22003"
)

func IsPqError(err error){
	var pqError *pq.Error

	if errors.As(err, &pqError){
		switch pqError.Code{

		}
	}
}


