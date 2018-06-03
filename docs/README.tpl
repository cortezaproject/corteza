{foreach $apis as $api}
# {api.title}

{api.description}

{foreach $api.apis as $name => $call}
	## {call.title}

	{call.description}

	### Request parameters

	| Parameter | Type | Method | Description | Default | Required? |
        | --------- | ---- | ------ | ----------- | ------- | --------- |
	{foreach $call.parameters as $method => $params}
		{foreach $params as $param}
		| {param.name} | {param.type} | {method|toupper} | {param.title} | {if empty($param.default)}N/A{else}{param.default}{/if}
{if $param.values}

Values:

{foreach $param.values as $value}- `{value}`
{/foreach}
{/if} | {if $param.required}YES{else}NO{/if} |
		{/foreach}
	{/foreach}


{/foreach}

{/foreach}