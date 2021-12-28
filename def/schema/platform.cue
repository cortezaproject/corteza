package schema

#platform: {
	ident: #baseHandle | *"corteza"

	components: [...{platform: ident} & #component]

	// env-var definitions
	// options: {}

	//

	// automation: {
	//  types: ....
	//  function ....
	// }
}
