The following block is literal (no template constructs work inside the tags)

<script type="text/template">
You can go {nuts} in here, no variables/etc will be parsed.
Except for {_CONSTANTS}, those will work.
</script>

Or this:

<script type="text/x-jquery">
You can go {nuts} in here, no variables/etc will be parsed.
{* And comments get removed also *}
</script>


But this works:

<script type="text/javascript">
// var some_object = {hello_this_is_a_variable};
</script>