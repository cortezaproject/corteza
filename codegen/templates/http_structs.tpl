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
	// {api.title} - {api.description}
	{struct.name} struct {
{foreach $struct.fields as $field}
		{field.name} {field.type} `{if $field.tag}{$field.tag} {/if}db:"{if $field.dbname}{$field.dbname}{else}{$field.name|decamel}{/if}"`{newline}
{/foreach}

		changed []string
	}

{/if}{/foreach}
)

{foreach $structs as $struct}
// New constructs a new instance of {struct.name}
func ({struct.name}) New() *{struct.name} {
	return &{struct.name}{}
}
{/foreach}

{foreach $structs as $struct}
{foreach $struct.fields as $field}
// Get the value of {field.name}
func ({self} *{struct.name}) Get{field.name}() {field.type} {
	return {self}.{field.name}
}

// Set the value of {field.name}
func ({self} *{struct.name}) Set{field.name}(value {field.type}) *{struct.name} {{if !$field.complex}
	if {self}.{field.name} != value {
		{self}.changed = append({self}.changed, "{field.name}")
		{self}.{field.name} = value
	}
{else}
	{self}.changed = append({self}.changed, "{field.name}")
	{self}.{field.name} = value
{/if}
	return {self}
}
{/foreach}

// Changes returns the names of changed fields
func ({self} *{struct.name}) Changes() []string {
	return {self}.changed
}

{/foreach}
