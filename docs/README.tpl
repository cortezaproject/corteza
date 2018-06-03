{foreach $apis as $api}
# {api.title}

{api.description}

{foreach $api.apis as $name => $call}
	## {call.title}

	{call.description}
{/foreach}

{/foreach}