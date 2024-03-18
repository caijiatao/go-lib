package main

type pointDB struct {
	value int
}

type A struct {
	db *pointDB
}

func (a *A) setValue(value int) {
	a.db.value = value
}

func (a *A) setPointDB(pdb *pointDB) {
	a.db = pdb
}

type B struct {
	db *pointDB
}

func main() {
	a := &A{
		db: &pointDB{
			value: 1,
		},
	}
	b := &B{
		db: a.db,
	}

	a.setValue(2)
	println(b.db.value)

	a.setPointDB(&pointDB{
		value: 3,
	})
	a.setValue(4)
	println(b.db.value)
}
