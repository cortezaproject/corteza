package {package}

{load warning_short.tpl}

{if !empty($imports)}
import (
{foreach ($imports as $import)}
       "{import}"
{/foreach}
)
{/if}

{if !empty($structs)}
type (
{foreach $structs as $struct}
       // {api.title} - {api.description}
       {struct.name} struct {
{foreach $struct.fields as $field}
               {field.name} {field.type} `{if $field.tag}{$field.tag} {/if}db:"{if $field.db}{$field.db}{else}{$field.name|decamel}{/if}"`{newline}
{/foreach}
       }

{/foreach}
)
{/if}
