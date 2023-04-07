package envoyx

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import ()

var (
	// needyResources is a list of resources that require a parent resource
	//
	// This list is primarily used when figuring out what nodes the dep. graph
	// should return when traversing.
	needyResources = map[string]bool{

		"corteza::compose:chart":        true,
		"corteza::compose:module":       true,
		"corteza::compose:module-field": true,

		"corteza::compose:page":        true,
		"corteza::compose:page-layout": true,
		"corteza::compose:record":      true,

		"corteza::compose:record-datasource": true,
	}

	// superNeedyResources is the second level of filtering in case the first
	// pass removes everything
	superNeedyResources = map[string]bool{
		"corteza::compose:module-field": true,
	}
)
