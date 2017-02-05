package db

func StoreCorrelationPair(base, assoc string) {
	execute(queries.StoreCorrelationPair, base, assoc)
}
