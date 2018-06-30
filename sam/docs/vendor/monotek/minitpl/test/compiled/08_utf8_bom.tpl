<?php $_v=&$this->vars;?>
This file contains utf8 BOM.

It also uses the <?php echo $_v['ldelim'];?>
load<?php echo $_v['rdelim'];?>
 construct:

<?php $this->push();$this->load("07_eval.tpl");$this->assign($_v);$this->render();$this->pop();?>
