package {package}

{load warning_short.tpl}

{if !empty($imports)}
import (
{foreach ($imports as $import)}
	{import}{EOL}{/foreach}
)
{/if}

{if !empty($structs)}
type (
{foreach $structs as $struct}
       // {api.title} - {api.description}
       {struct.name} struct {
{foreach $struct.fields as $field}
               {field.name} {field.type} `json:"{if $field.json}{$field.json}{elseif $field.db}{$field.db}{else}{$field.name|decamel}{/if}{if $field.omitempty},omitempty{/if}" db:"{if $field.db}{$field.db}{else}{$field.name|decamel}{/if}"`{newline}
{/foreach}
       }

{/foreach}
)
{/if}
