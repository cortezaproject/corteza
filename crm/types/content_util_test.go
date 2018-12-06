package types

import (
	"testing"
)

func TestContentReport(t *testing.T) {
	r := &ContentReport{}
	r.ScanMetrics("alias:exp")
	if len(r.Metrics) == 0 {
		t.Log("No metrics scanned")
		t.Fail()
	} else if r.Metrics[0].Alias != "alias" {
		t.Log("Alias not parsed")
		t.Fail()
	} else if r.Metrics[0].Expression != "exp" {
		t.Log("Expression not parsed")
		t.Fail()
	}

	r.ScanMetrics("exp")
	if len(r.Metrics) == 0 {
		t.Log("No metrics scanned")
		t.Fail()
	} else if r.Metrics[0].Alias != "" {
		t.Log("Alias should be empty")
		t.Fail()
	} else if r.Metrics[0].Expression != "exp" {
		t.Log("Expression not parsed")
		t.Fail()
	}

	r.ScanDimensions("alias:field|m1|m2")
	if len(r.Dimensions) == 0 {
		t.Log("No dimensions scanned")
		t.Fail()
	} else if r.Dimensions[0].Alias != "alias" {
		t.Log("Alias not parsed")
		t.Fail()
	} else if r.Dimensions[0].Field != "field" {
		t.Log("Expression not parsed")
		t.Fail()
	} else if len(r.Dimensions[0].Modifiers) != 2 && r.Dimensions[0].Modifiers[0] != "m1" || r.Dimensions[0].Modifiers[1] != "m2" {
		t.Log("Modifiers not parsed")
		t.Fail()
	}

	r.ScanDimensions("field")
	if len(r.Dimensions) == 0 {
		t.Log("No dimensions scanned")
		t.Fail()
	} else if r.Dimensions[0].Alias != "" {
		t.Log("Alias should be empty")
		t.Fail()
	} else if len(r.Dimensions[0].Modifiers) > 0 {
		t.Log("Modifiers should be empty")
		t.Fail()
	}
}
