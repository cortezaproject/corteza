<?php $_v=&$this->vars; $_v['key']="nu";?>




<?php $_v['i']=1;if(!empty($_v['items']))foreach($_v['items'] as $_v['index'] => $_v['item']){?>
	<?php $_v['id'] = $_v['key'] . $_v['i'];?>

	<?php if($_v['i']>1){ echo $_v['i'];?>
nd<?php }else{?>
first<?php }?>
	<?php $_v['i']++;}?>
