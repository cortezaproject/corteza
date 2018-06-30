<?php $_v=&$this->vars;?>
The following block is literal (no template constructs work inside the tags)

<script type="text/template">
You can go {nuts} in here, no variables/etc will be parsed.
Except for <?php echo _CONSTANTS;?>
, those will work.
</script>

Or this:

<script type="text/x-jquery">
You can go {nuts} in here, no variables/etc will be parsed.

</script>


But this works:

<script type="text/javascript">
// var some_object = <?php echo $_v['hello_this_is_a_variable'];?>
;
</script>