#### 1.3. Modifiers

Modifiers are PHP functions, which take one arbitrary value
(usually, string), and return a string, which gets shown.
Some functions in PHP can be used as modifiers
(strtoupper, ucfirst, str_rot13, strrev, count, ...).

~~~~~~~~~~~~~~~~~~~~~~~~
/** Custom modifier example */

function add_it_up($array)
{
	$size = 0;
	foreach ($array as $value) {
		$size += $value['size'];
	}
	return $size;
}
~~~~~~~~~~~~~~~~~~~~~~~~

To use the above function in a template, just do `{$variable|add_it_up}`,
which will loop trough all `$variable` items and add up the value of
the size element and return the total sum of all sizes.

~~~~~~~~~~~~~~
{variable|escape}
{variable|toupper}
{variable|add_id_up}
~~~~~~~~~~~~~~

The above code gets compiled to:

~~~~~~~~~~~~~~
<?php
	$_v = &$this->vars;
	echo htmlspecialchars($_v['variable'], ENT_QUOTES);
	echo strtoupper($_v['variable']);
	echo add_id_up($_v['variable']);
?>
~~~~~~~~~~~~~~

A common use for modifiers is outputting data for javascript,
using the php function `json_encode`.
