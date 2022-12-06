<template>
  <nylas-email
    :id="componentID"
    :access_token="accessToken"
    theme="light"
  />
</template>

<script>
import base from './base'
import '@nylas/components-email'

export default {
  extends: base,
}
</script>
