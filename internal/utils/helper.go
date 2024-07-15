package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func CommitOrRollback(ctx *fiber.Ctx, tx pgx.Tx) {
	err := recover()
	if err != nil {
		_ = tx.Rollback(ctx.Context())
		//PanicIfError(ctx, fiber.StatusInternalServerError, errorRollback)
		//panic(errorRollback)
	} else {
		_ = tx.Commit(ctx.Context())
		//PanicIfError(ctx, fiber.StatusInternalServerError, errorCommit)
		//panic(errorCommit)
	}
}
