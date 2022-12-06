<template>
  <nylas-conversation
    :id="componentID"
    :access_token="accessToken"
    theme="light"
  />
</template>

<script>
import base from './base'
import '@nylas/components-conversation'

export default {
  extends: base,
}
</script>
