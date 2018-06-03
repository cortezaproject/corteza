#### 1.5. Constants

Value starting with `_` is assumed to be a constant. The template
system will output the value of the constant if defined,
or just the name of the constant. Item `{_MY_CONSTANT}` will be
used as a constant, because of those rules, however `{$_my_variable}`
wouldnt be, since it starts with the variable identifier `$`.

~~~~~~~~~~~
{_MY_CONSTANT}
{_this_is_also_a_constant}
{$_my_variable}
~~~~~~~~~~~