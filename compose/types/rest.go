package types

type (
	BulkRecord struct {
		RefField string    `json:"refField,omitempty"`
		Set      RecordSet `json:"set,omitempty"`
	}

	BulkRecordSet []*BulkRecord
)

func (set BulkRecordSet) ToBulkOperations(dftModule uint64, dftNamespace uint64) (oo []*BulkRecordOperation, err error) {
	for _, br := range set {
		for _, rr := range br.Set {
			// No use in allowing cross-namespace record creation.
			rr.NamespaceID = dftNamespace

			// default module
			if rr.ModuleID == 0 {
				rr.ModuleID = dftModule
			}
			b := &BulkRecordOperation{
				Record:    rr,
				Operation: OperationTypeUpdate,
				LinkBy:    br.RefField,
			}

			// If no RecordID is defined, we should create it
			if rr.ID == 0 {
				b.Operation = OperationTypeCreate
			}

			if rr.DeletedAt != nil {
				b.Operation = OperationTypeDelete
			}

			oo = append(oo, b)
		}
	}

	return
}
