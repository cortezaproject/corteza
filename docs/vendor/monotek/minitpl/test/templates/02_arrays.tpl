#### 1.2. Arrays

There is a shorthand syntax for using arrays inside a template.
The array index separator is a dot (`.`). Consecutive dots can be
used for traversing into array depth like `{$array.items.0}`, or
even variables starting with `$`, like `{$sections.$news_section.title}`.
PHP syntax for arrays is also supported.

~~~~~~~~~~~~
{$array.items} is the same as {$array['items']}
{$array.items.0} is the same as {$array['items'][0]}
{$array.$items.0} is the same as {$array[$items][0]}
~~~~~~~~~~~~
