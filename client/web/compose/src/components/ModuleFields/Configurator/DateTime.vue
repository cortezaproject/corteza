<template>
  <div>
    <div>
      <b-form-checkbox
        v-model="f.options.onlyDate"
        :disabled="f.options.onlyTime"
      >
        {{ $t('kind.dateTime.dateOnly') }}
      </b-form-checkbox>
    </div>
    <div>
      <b-form-checkbox
        v-model="f.options.onlyTime"
        :disabled="f.options.onlyDate"
      >
        {{ $t('kind.dateTime.timeOnly') }}
      </b-form-checkbox>
    </div>
    <div>
      <b-form-checkbox
        v-model="f.options.onlyPastValues"
        :disabled="f.options.onlyFutureValues"
      >
        {{ $t('kind.dateTime.pastValuesOnly') }}
      </b-form-checkbox>
    </div>
    <div>
      <b-form-checkbox
        v-model="f.options.onlyFutureValues"
        :disabled="f.options.onlyPastValues"
      >
        {{ $t('kind.dateTime.futureValuesOnly') }}
      </b-form-checkbox>
    </div>
    <div>
      <b-form-checkbox v-model="f.options.outputRelative">
        {{ $t('kind.dateTime.relativeOutput') }}
      </b-form-checkbox>
    </div>
    <div>
      {{ $t('kind.dateTime.outputFormat') }}<b-form-input
        v-model="f.options.format"
        plain
        placeholder="YYYY-MM-DD HH:ii"
      />
    </div>

    <i18next
      path="kind.dateTime.outputFormatFootnote"
      tag="label"
    >
      <a
        href="https://momentjs.com/docs/#/displaying/format/"
        target="_blank"
        rel="noopener noreferrer"
      >Moment.js</a>
    </i18next>
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
