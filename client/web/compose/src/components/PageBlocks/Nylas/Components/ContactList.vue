<template>
  <nylas-contact-list
    :id="componentID"
    :access_token="accessToken"
    theme="light"
  />
</template>

<script>
import base from './base'
import '@nylas/components-contact-list'

export default {
  extends: base,
}
</script>
