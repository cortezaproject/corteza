<?php $_v=&$this->vars;?>
#### 1.2. Arrays

There is a shorthand syntax for using arrays inside a template.
The array index separator is a dot (`.`). Consecutive dots can be
used for traversing into array depth like `<?php echo $_v['array']['items']['0'];?>
`, or
even variables starting with `$`, like `<?php echo $_v['sections'][$_v['news_section']]['title'];?>
`.
PHP syntax for arrays is also supported.

~~~~~~~~~~~~
<?php echo $_v['array']['items'];?>
 is the same as <?php echo $_v['array']['items'];echo $_v['array']['items']['0'];?>
 is the same as <?php echo $_v['array']['items'][0];echo $_v['array'][$_v['items']]['0'];?>
 is the same as <?php echo $_v['array'][$_v['items']][0];?>

~~~~~~~~~~~~
