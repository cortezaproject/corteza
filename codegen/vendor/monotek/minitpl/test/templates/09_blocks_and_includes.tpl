{inline snippet1}
This is like a copy-paste over your {eval echo __FILE__}, anywhere you use this construct.
{/inline}

Functions bro, functions?

{eval $i=1}
{block function_call}
This is a created function. This is the {i++}st/nd/rd/th time calling this method.
{if $i <= 5}
{block:function_call}
{/if}
{/block}

{inline:snippet1}

{inline:snippet1}

{block:function_call}