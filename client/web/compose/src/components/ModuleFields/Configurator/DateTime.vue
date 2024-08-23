<template>
  <div>
    <b-form-group
      :label="$t('kind.dateTime.type.label')"
      label-class="text-primary"
      class="w-25"
    >
      <b-form-radio-group
        v-b-tooltip.hover="{ title: hasData ? $t('not-configurable') : '', placement: 'left', container: '#body' }"
        :checked="inputType"
        :options="[
          { value: 'dateTime', text: $t('kind.dateTime.type.options.dateTime') },
          { value: 'date', text: $t('kind.dateTime.type.options.date') },
          { value: 'time', text: $t('kind.dateTime.type.options.time') },
        ]"
        :disabled="hasData"
        stacked
        @input="onTypeChange"
      />
    </b-form-group>

    <b-form-group
      :label="$t('kind.dateTime.constraints.label')"
      label-class="text-primary"
      class="mt-2"
    >
      <b-form-radio-group
        :checked="constraintType"
        :options="[
          { value: 'all', text: $t('kind.dateTime.constraints.options.all') },
          { value: 'pastValuesOnly', text: $t('kind.dateTime.constraints.options.pastValuesOnly') },
          { value: 'futureValuesOnly', text: $t('kind.dateTime.constraints.options.futureValuesOnly') },
        ]"
        stacked
        @input="onConstraintChange"
      />
    </b-form-group>

    <b-form-group
      :label="$t('kind.dateTime.outputFormat')"
      label-class="text-primary"
      class="mt-2"
    >
      <b-form-checkbox
        v-model="f.options.outputRelative"
        class="mb-2"
      >
        {{ $t('kind.dateTime.relativeOutput') }}
      </b-form-checkbox>

      <template v-if="!f.options.outputRelative">
        <b-form-input
          v-model="f.options.format"
          plain
          placeholder="YYYY-MM-DD HH:ii"
        />

        <div class="small text-muted">
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
    </b-form-group>
  </div>
</template>

<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,

  computed: {
    inputType () {
      if (this.f.options.onlyDate) return 'date'
      if (this.f.options.onlyTime) return 'time'

      return 'dateTime'
    },

    constraintType () {
      if (this.f.options.onlyPastValues) return 'pastValuesOnly'
      if (this.f.options.onlyFutureValues) return 'futureValuesOnly'

      return 'all'
    },
  },

  methods: {
    onTypeChange (v) {
      this.f.options.onlyDate = v === 'date'
      this.f.options.onlyTime = v === 'time'
    },

    onConstraintChange (v) {
      this.f.options.onlyPastValues = v === 'pastValuesOnly'
      this.f.options.onlyFutureValues = v === 'futureValuesOnly'
    },
  },
}
</script>
