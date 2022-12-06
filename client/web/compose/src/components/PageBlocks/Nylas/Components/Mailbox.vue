<template>
  <nylas-mailbox
    :id="componentID"
    :access_token="accessToken"
    theme="light"
  />
</template>

<script>
import base from './base'
import '@nylas/components-mailbox'

export default {
  extends: base,
}
</script>
