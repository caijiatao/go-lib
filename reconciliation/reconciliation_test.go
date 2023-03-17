package reconciliation

import (
	"log"
	"testing"
)

func TestReconciliation(t *testing.T) {
	ron := NewReconciliationImpl()

	ron.RegisterReconciliation(ReconciliationConfig{})

	err := ron.Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("reconciliation run end")
}
