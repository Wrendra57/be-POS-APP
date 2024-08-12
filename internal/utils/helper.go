package utils

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func CommitOrRollback(ctx context.Context, tx pgx.Tx) {
	err := recover()
	if err != nil {
		errRoleback := tx.Rollback(ctx)
		PanicIfError(errRoleback)
		panic(err)
	} else {
		errCommit := tx.Commit(ctx)
		PanicIfError(errCommit)
	}
}
