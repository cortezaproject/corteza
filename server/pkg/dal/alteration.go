package dal

import "encoding/json"

type (
	Alteration struct {
		ID           uint64
		BatchID      uint64
		DependsOn    uint64
		Resource     string
		ResourceType string
		ConnectionID uint64

		AttributeAdd      *AttributeAdd
		AttributeDelete   *AttributeDelete
		AttributeReType   *AttributeReType
		AttributeReEncode *AttributeReEncode
		ModelAdd          *ModelAdd
		ModelDelete       *ModelDelete
	}

	AlterationSet []*Alteration

	AttributeAdd struct {
		Attr *Attribute `json:"attr"`
	}

	AttributeDelete struct {
		Attr *Attribute `json:"attr"`
	}

	AttributeReType struct {
		Attr *Attribute `json:"attr"`
		To   Type       `json:"to"`
	}

	AttributeReEncode struct {
		Attr *Attribute `json:"attr"`
		To   *Attribute `json:"to"`
	}

	ModelAdd struct {
		Model *Model `json:"model"`
	}

	ModelDelete struct {
		Model *Model `json:"model"`
	}

	// auxAttributeType is a helper struct used for marshaling/unmarshaling
	//
	// This is required since the Type inside the Attribute is an interface and we
	// need to help the encoding/json a bit.
	auxAttributeReType struct {
		Attr *Attribute
		To   *auxAttributeType
	}
)

// Merge merges the two alteration slices
func (a AlterationSet) Merge(b AlterationSet) (c AlterationSet) {
	// @todo don't blindly append the two slices since there can be duplicates
	// or overlapping alterations which would cause needles processing
	//
	// A quick list of overlapping alterations:
	// * attribute A added and then renamed from A to A'
	// * attribute A renamed to A' and then renamed to A''
	// * attribute A deleted and then created
	//
	// For now we'll simply append them and worry about improvements on a later stage

	return append(a, b...)
}

func (a AttributeReType) MarshalJSON() ([]byte, error) {
	aux := auxAttributeReType{
		Attr: a.Attr,
		To:   &auxAttributeType{},
	}

	switch t := a.To.(type) {
	case *TypeID:
		aux.To.Type = "ID"
		aux.To.ID = t

	case *TypeRef:
		aux.To.Type = "Ref"
		aux.To.Ref = t

	case *TypeTimestamp:
		aux.To.Type = "Timestamp"
		aux.To.Timestamp = t

	case *TypeTime:
		aux.To.Type = "Time"
		aux.To.Time = t

	case *TypeDate:
		aux.To.Type = "Date"
		aux.To.Date = t

	case *TypeNumber:
		aux.To.Type = "Number"
		aux.To.Number = t

	case *TypeText:
		aux.To.Type = "Text"
		aux.To.Text = t

	case *TypeBoolean:
		aux.To.Type = "Boolean"
		aux.To.Boolean = t

	case *TypeEnum:
		aux.To.Type = "Enum"
		aux.To.Enum = t

	case *TypeGeometry:
		aux.To.Type = "Geometry"
		aux.To.Geometry = t

	case *TypeJSON:
		aux.To.Type = "JSON"
		aux.To.JSON = t

	case *TypeBlob:
		aux.To.Type = "Blob"
		aux.To.Blob = t

	case *TypeUUID:
		aux.To.Type = "UUID"
		aux.To.UUID = t
	}

	return json.Marshal(aux)
}

func (a *AttributeReType) UnmarshalJSON(data []byte) (err error) {
	aux := &auxAttributeReType{}
	err = json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if a == nil {
		*a = AttributeReType{}
	}

	a.Attr = aux.Attr

	switch aux.To.Type {
	case "ID":
		a.To = aux.To.ID

	case "Ref":
		a.To = aux.To.Ref

	case "Timestamp":
		a.To = aux.To.Timestamp

	case "Time":
		a.To = aux.To.Time

	case "Date":
		a.To = aux.To.Date

	case "Number":
		a.To = aux.To.Number

	case "Text":
		a.To = aux.To.Text

	case "Boolean":
		a.To = aux.To.Boolean

	case "Enum":
		a.To = aux.To.Enum

	case "Geometry":
		a.To = aux.To.Geometry

	case "JSON":
		a.To = aux.To.JSON

	case "Blob":
		a.To = aux.To.Blob

	case "UUID":
		a.To = aux.To.UUID
	}

	return
}
