package {package}

{if strpos($name, "Literal") !== false}
	{foreach $fields as $field}
		{field}{newline}
	{/foreach}
{else}
	// {api.title}
	type {name} struct {
{foreach $fields as $field}
		{field.name} {field.type}{if $field.tag} `{$field.tag}`{/if}{newline}
{/foreach}
	}

func ({name}) new() *{name} {
	return &{name}{}
}

{/if}

{foreach $fields as $field}
func ({self} *{name}) Get{field.name}() {field.type} {
	return {self}.{field.name}
}

func ({self} *{name}) Set{field.name}(value {field.type}) *{name} {
	{self}.{field.name} = value
	return {self}
}
{/foreach}
