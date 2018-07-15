{for $i=99; $i>=1; $i--}
{i} bottles of beer on the wall, {i} bottles of beer. Take one down and pass it around, {eval echo $i-1} bottles of beer on the wall.
{/for}

{if $i > 10}
	We have some beer left.
{elseif $i > 3}
	Our beer is going to run out.
{else}
	Critical! Only {i} beer left.
{/if}