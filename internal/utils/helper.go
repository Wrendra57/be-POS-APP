package utils

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func CommitOrRollback(ctx context.Context, tx pgx.Tx) {
	err := recover()
	if err != nil {
		fmt.Println("roleback")
		errRoleback := tx.Rollback(ctx)
		PanicIfError(errRoleback)
	} else {
		errCommit := tx.Commit(ctx)
		PanicIfError(errCommit)
	}
}
