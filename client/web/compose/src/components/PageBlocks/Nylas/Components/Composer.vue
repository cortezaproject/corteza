<template>
  <nylas-composer
    :id="componentID"
    :access_token="accessToken"
    theme="light"
  />
</template>

<script>
import base from './base'
import '@nylas/components-composer'

export default {
  extends: base,
}
</script>
