{eval $key="nu"}

{inline item}
{if $i>1}{i}nd{else}first{/if}
{/inline}

{eval $i=1}
{foreach $items as $index => $item}
	{eval $id = $key . $i}
	{inline:item}
	{eval $i++}
{/foreach}
