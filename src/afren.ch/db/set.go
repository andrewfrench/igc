package db

func AssociationBaseExists(base string) bool {
	row := queryRow(queries.GetAssociationCount, base)

	var count int
	row.Scan(&count)

	return count > 0
}

func InsertAssociationSet(base string, data []byte) {
	execute(queries.InsertAssociationSet, base, data)
}

func UpdateAssociationSet(base string, data []byte ){
	execute(queries.UpdateAssociationSet, base, data)
}

func QueryIncomingSet() (string, []string) {
	rows := queryRows(queries.QueryIncomingSet)

	var base string
	var assocs []string

	for rows.Next() {
		var assoc string
		rows.Scan(&base, &assoc)
		assocs = append(assocs, assoc)
	}

	return base, assocs
}

func QueryAssociationSet(base string) []byte {
	row := queryRow(queries.QueryAssociationSet, base)

	var assoc []byte
	row.Scan(&assoc)

	return assoc
}