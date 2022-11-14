<template>
  <p>
    {{ $t('editor.unsupported') }}
  </p>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'system.templates',
    keyPrefix: 'editor.content.editor',
  },
}
</script>
