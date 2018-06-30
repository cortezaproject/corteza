There is a distinct difference between
calling `$tpl->get` or `$tplx->get`.

When `$tpl` is a global variable, the template engine detects this.

{$tpl->_paths} {$tplx->_paths};