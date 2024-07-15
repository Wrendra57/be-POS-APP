package config

import "github.com/jackc/pgx/v5"

func TxConfig() pgx.TxOptions {
	txOptions := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}
	return txOptions
}
