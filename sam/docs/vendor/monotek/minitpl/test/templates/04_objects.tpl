#### 1.4. Objects

You can also call object methods as modifiers, from global or
local objects. The compiler will determine at run-time, if
you have a global object by the name, and modify the code
accordingly.

~~~~~~~~~~~
{$variable|$memcache->get}
~~~~~~~~~~~

If the global object doesn't exist, it assumes a local object.

~~~~~~~~~~~
<?php
	$_v = &$this->vars;
	echo $_v['memcache']->get($_v['variable']);
?>
~~~~~~~~~~~

However, if a global object by the name $memcache exists at
compile time:

~~~~~~~~~~
<?php
	$_v = &$this->vars;
	global $memcache;
	echo $memcache->get($_v['variable']);
?>
~~~~~~~~~~

You can also use variables from objects, in the same way.

~~~~~~~~~~~~
{$memcache->variable}
~~~~~~~~~~~~

Compiles to one of theese:

~~~~~~~~~~
<?php
	$_v = &$this->vars;
	echo $_v['memcache']->variable;
?>
~~~~~~~~~~

~~~~~~~~~~
<?php
	$_v = &$this->vars;
	global $memcache;
	echo $memcache->variable;
?>
~~~~~~~~~~
