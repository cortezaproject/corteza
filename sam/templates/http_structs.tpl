package {package}

{load warning.tpl}

{if !empty($imports)}
import (
{foreach ($imports as $import)}
	"{import}"
{/foreach}
)
{/if}

type ({foreach $structs as $struct}

{if strpos($name, "Literal") !== false}
	{foreach $struct.fields as $field}
		{field}{newline}
	{/foreach}
{else}
	// {api.title}
	{struct.name} struct {
{foreach $struct.fields as $field}
		{field.name} {field.type}{if $field.tag || $field.dbname} `{$field.tag}{if $field.dbname} db:"{$field.dbname}"{/if}`{/if}{newline}
{/foreach}

		changed []string
	}

{/if}{/foreach}
)

/* Constructors */
{foreach $structs as $struct}
func ({struct.name}) New() *{struct.name} {
	return &{struct.name}{}
}
{/foreach}

/* Getters/setters */
{foreach $structs as $struct}
{foreach $struct.fields as $field}
func ({self} *{struct.name}) Get{field.name}() {field.type} {
	return {self}.{field.name}
}

func ({self} *{struct.name}) Set{field.name}(value {field.type}) *{struct.name} {{if !$field.complex}
	if {self}.{field.name} != value {
		{self}.changed = append({self}.changed, "{field.name}")
		{self}.{field.name} = value
	}
{else}
	{self}.{field.name} = value
{/if}
	return {self}
}
{/foreach}

{/foreach}
