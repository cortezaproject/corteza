<?php $_v=&$this->vars;?>
We check if arrays are empty to suppress warnings/notices.

<?php foreach(array("GET", "POST") as $_v['method']){?>
HTTP <?php echo $_v['method'];}foreach(array("GET", "POST") as $_v['method']){?>
HTTP <?php echo $_v['method'];}if(!empty($_v['methods']))foreach($_v['methods'] as $_v['method']){?>
HTTP <?php echo $_v['method'];}?>
