<?php $_v=&$this->vars;?>
#### 1.1. Variables

Using variables from templates is easy. Variables are enclosed
in curly braces like so: `<?php echo $_v['variable'];?>
`. Since variables are parsed
on the last stage, prefixing variables with the dollar sign is
optional. So, `<?php echo $_v['variable'];?>
` is the same as `<?php echo $_v['variable'];?>
`.

~~~~~~~~~~~~~
<?php echo $_v['variable'];?>
 is the same as <?php echo $_v['variable'];?>

~~~~~~~~~~~~~
