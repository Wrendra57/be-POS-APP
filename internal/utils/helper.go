package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func CommitOrRollback(ctx *fiber.Ctx, tx pgx.Tx) {
	err := recover()
	if err != nil {
		fmt.Println("roleback")
		errRoleback := tx.Rollback(ctx.Context())
		PanicIfError(errRoleback)
	} else {
		errCommit := tx.Commit(ctx.Context())
		PanicIfError(errCommit)
	}
}
