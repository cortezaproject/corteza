<template>
  <div>
    <div>
      <b-form-checkbox v-model="f.options.trimPath">
        {{ $t('kind.url.label') }}
      </b-form-checkbox>
    </div>
    <div>
      <b-form-checkbox v-model="f.options.trimFragment">
        {{ $t('kind.url.trimHash') }}
      </b-form-checkbox>
    </div>
    <div>
      <b-form-checkbox v-model="f.options.trimQuery">
        {{ $t('kind.url.trimQuestionMark') }}
      </b-form-checkbox>
    </div>
    <div>
      <b-form-checkbox v-model="f.options.onlySecure">
        {{ $t('kind.url.sshOnly') }}
      </b-form-checkbox>
    </div>
    <div>
      <b-form-checkbox v-model="f.options.outputPlain">
        {{ $t('kind.url.preventToLink') }}
      </b-form-checkbox>
    </div>
  </div>
</template>

<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,
}
</script>
