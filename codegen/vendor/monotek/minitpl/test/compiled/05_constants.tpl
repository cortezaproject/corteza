<?php $_v=&$this->vars;?>
#### 1.5. Constants

Value starting with `_` is assumed to be a constant. The template
system will output the value of the constant if defined,
or just the name of the constant. Item `<?php echo _MY_CONSTANT;?>
` will be
used as a constant, because of those rules, however `<?php echo $_v['_my_variable'];?>
`
wouldnt be, since it starts with the variable identifier `$`.

~~~~~~~~~~~
<?php echo _MY_CONSTANT;echo _this_is_also_a_constant;echo $_v['_my_variable'];?>

~~~~~~~~~~~