<?php $_v=&$this->vars; global $tpl;?>
There is a distinct difference between
calling `$tpl->get` or `$tplx->get`.

When `$tpl` is a global variable, the template engine detects this.

<?php echo $tpl->_paths;?>
 <?php echo $_v['tplx']->_paths;?>
;