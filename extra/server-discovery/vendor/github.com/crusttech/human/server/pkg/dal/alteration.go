package dal

import (
	"encoding/json"
	"fmt"
)

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
		To   Codec      `json:"to"`
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
		Attr *Attribute        `json:"attr"`
		To   *auxAttributeType `json:"to"`
	}

	// auxAttributeReEncode is a helper struct used for marshaling/unmarshaling
	//
	// This is required since the Codec inside the Attribute is an interface and we
	// need to help the encoding/json a bit.
	auxAttributeReEncode struct {
		Attr *Attribute     `json:"attr"`
		To   *auxStoreCodec `json:"to"`
	}
)

// Merge merges the two alteration slices
func (aa AlterationSet) Merge(bb AlterationSet) (cc AlterationSet) {
	// @todo the Merge function currently just merges the two slices together
	//       and removes any duplicates that would occur due to the merge.
	//       We should also handle overlapping/transitive alterations to reduce
	//       the amount of needles processing.
	//
	// A quick list of overlapping alterations:
	// * attribute A added and then renamed from A to A'
	// * attribute A renamed to A' and then renamed to A''
	// * attribute A deleted and then created

	aux := append(aa, bb...)
	seen := make(map[int]bool, len(aux)/2)

	for i, a := range aux {
		if seen[i] {
			continue
		}

		found := false
		for j := i + 1; j < len(aux); j++ {
			if seen[j] {
				continue
			}

			if a.compare(*aux[j]) {
				if !found {
					cc = append(cc, aux[j])
				}

				seen[j] = true
				found = true
			}
		}

		if !found {
			cc = append(cc, a)
		}
	}

	return
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

func (a AttributeReEncode) MarshalJSON() ([]byte, error) {
	aux := auxAttributeReEncode{
		Attr: a.Attr,
		To:   &auxStoreCodec{},
	}

	switch t := a.To.(type) {
	case *CodecPlain:
		aux.To.Type = "CodecPlain"
		aux.To.CodecPlain = t
	case *CodecRecordValueSetJSON:
		aux.To.Type = "CodecRecordValueSetJSON"
		aux.To.CodecRecordValueSetJSON = t
	case *CodecAlias:
		aux.To.Type = "CodecAlias"
		aux.To.CodecAlias = t
	}

	return json.Marshal(aux)
}

func (a *AttributeReEncode) UnmarshalJSON(data []byte) (err error) {
	aux := &auxAttributeReEncode{}
	err = json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if a == nil {
		*a = AttributeReEncode{}
	}

	a.Attr = aux.Attr

	switch aux.To.Type {
	case "CodecPlain":
		a.To = aux.To.CodecPlain

	case "CodecRecordValueSetJSON":
		a.To = aux.To.CodecRecordValueSetJSON

	case "CodecAlias":
		a.To = aux.To.CodecAlias
	}

	return
}

func (a Alteration) compare(b Alteration) (cmp bool) {
	if a.AttributeAdd == nil && b.AttributeAdd != nil {
		return false
	}
	if a.AttributeDelete == nil && b.AttributeDelete != nil {
		return false
	}
	if a.AttributeReType == nil && b.AttributeReType != nil {
		return false
	}
	if a.AttributeReEncode == nil && b.AttributeReEncode != nil {
		return false
	}
	if a.ModelAdd == nil && b.ModelAdd != nil {
		return false
	}
	if a.ModelDelete == nil && b.ModelDelete != nil {
		return false
	}

	switch {
	case a.AttributeAdd != nil:
		return a.compareAttributeAdd(b)
	case a.AttributeDelete != nil:
		return a.compareAttributeDelete(b)
	case a.AttributeReType != nil:
		return a.compareAttributeReType(b)
	case a.AttributeReEncode != nil:
		return a.compareAttributeReEncode(b)
	case a.ModelAdd != nil:
		return a.compareModelAdd(b)
	case a.ModelDelete != nil:
		return a.compareModelDelete(b)
	}

	panic(fmt.Sprintf("unsupported alteration type %v", a))
}

func (a Alteration) compareAttributeAdd(b Alteration) bool {
	if a.AttributeAdd == nil || b.AttributeAdd == nil {
		return a.AttributeAdd == b.AttributeAdd
	}
	return a.AttributeAdd.Attr.Compare(b.AttributeAdd.Attr)
}

func (a Alteration) compareAttributeDelete(b Alteration) bool {
	if a.AttributeDelete == nil || b.AttributeDelete == nil {
		return a.AttributeDelete == b.AttributeDelete
	}
	return a.AttributeDelete.Attr.Compare(b.AttributeDelete.Attr)
}

func (a Alteration) compareAttributeReType(b Alteration) bool {
	if !a.AttributeReType.Attr.Compare(b.AttributeReType.Attr) {
		return false
	}

	if a.AttributeReType == nil || b.AttributeReType == nil {
		return a.AttributeReType == b.AttributeReType
	}
	return a.AttributeReType.To.Type() == b.AttributeReType.To.Type()
}

func (a Alteration) compareAttributeReEncode(b Alteration) bool {
	if !a.AttributeReEncode.Attr.Compare(b.AttributeReEncode.Attr) {
		return false
	}

	if a.AttributeReEncode == nil || b.AttributeReEncode == nil {
		return a.AttributeReEncode == b.AttributeReEncode
	}
	return a.AttributeReEncode.To.Type() == b.AttributeReEncode.To.Type()
}

func (a Alteration) compareModelAdd(b Alteration) bool {
	if a.ModelAdd == nil || b.ModelAdd == nil {
		return a.ModelAdd == b.ModelAdd
	}
	return a.ModelAdd.Model.Compare(*b.ModelAdd.Model)
}

func (a Alteration) compareModelDelete(b Alteration) bool {
	if a.ModelDelete == nil || b.ModelDelete == nil {
		return a.ModelDelete == b.ModelDelete
	}
	return a.ModelDelete.Model.Compare(*b.ModelDelete.Model)
}
