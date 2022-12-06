<template>
  <nylas-agenda
    :id="componentID"
    :access_token="accessToken"
    theme="light"
  />
</template>

<script>
import base from './base'
import '@nylas/components-agenda'

export default {
  extends: base,
}
</script>
