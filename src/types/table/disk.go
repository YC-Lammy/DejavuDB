package table

func (t *Table) ToDisk(path string) {
	/*
		one file per column
		create one file every a 10k row
	*/
	for k, v := range t.Columns {
		go v.ToDisk(path)
	}
}
