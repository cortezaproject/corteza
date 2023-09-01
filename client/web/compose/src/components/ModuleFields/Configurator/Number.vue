<template>
  <div>
    <b-row>
      <b-col
        cols="12"
        sm="6"
      >
        <b-form-group
          :label="$t('kind.number.displayType.label')"
          label-class="text-primary"
        >
          <b-form-radio-group
            v-model="f.options.display"
            button-variant="outline-primary"
            :options="displayOptions"
            buttons
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        sm="6"
      >
        <b-form-group
          label-class="mb-2 text-primary"
        >
          <template #label>
            {{ `${$t('kind.number.precisionLabel')} ${(f.options.precision)}` }}
          </template>
          <b-form-input
            v-model="f.options.precision"
            :placeholder="$t('kind.number.precisionPlaceholder')"
            type="range"
            min="0"
            max="6"
            size="lg"
            class="mt-1 mb-2"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <hr>

    <b-row>
      <template v-if="f.options.display === 'number'">
        <b-col
          cols="12"
          sm="6"
        >
          <b-form-group
            :label="$t('kind.number.prefixLabel')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.prefix"
              :placeholder="$t('kind.number.prefixPlaceholder')"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          sm="6"
        >
          <b-form-group
            :label="$t('kind.number.suffixLabel')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.suffix"
              :placeholder="$t('kind.number.suffixPlaceholder')"
            />
          </b-form-group>
        </b-col>
      </template>

      <template v-if="f.options.display === 'number'">
        <b-col
          cols="12"
          sm="6"
        >
          <b-form-group
            :label="$t('kind.number.formatLabel')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.format"
              :placeholder="$t('kind.number.formatPlaceholder')"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
        >
          <b-form-group
            v-if="f.options.display === 'number'"
            :label="$t('kind.number.examplesLabel')"
            label-class="text-primary"
          >
            <table
              style="width: 100%;"
            >
              <tr>
                <th>{{ $t('kind.number.exampleInput') }}</th>
                <th>{{ $t('kind.number.exampleFormat') }}</th>
                <th>{{ $t('kind.number.exampleResult') }}</th>
              </tr>

              <tr>
                <td>10000.234</td>
                <td>$0.00</td>
                <td>$10000.23</td>
              </tr>

              <tr>
                <td>0.974878234</td>
                <td>0.000%</td>
                <td>97.488%</td>
              </tr>

              <tr>
                <td>1</td>
                <td>0%</td>
                <td>100%</td>
              </tr>

              <tr>
                <td>238</td>
                <td>00:00:00</td>
                <td>0:03:58</td>
              </tr>
            </table>
          </b-form-group>
        </b-col>
      </template>

      <template v-if="f.options.display === 'progress'">
        <b-col
          cols="12"
          sm="6"
        >
          <b-form-group
            :label="$t('kind.number.progress.minimumValue')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.min"
              type="number"
              number
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          sm="6"
        >
          <b-form-group
            :label="$t('kind.number.progress.maximumValue')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.max"
              type="number"
              number
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          sm="6"
        >
          <b-form-group
            :label="$t('kind.number.progress.variants.default')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="f.options.variant"
              :options="variants"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          sm="3"
        >
          <b-form-group
            class="mb-0"
          >
            <b-form-checkbox
              v-model="f.options.showValue"
              class="mb-2"
            >
              {{ $t('kind.number.progress.show.value') }}
            </b-form-checkbox>

            <b-form-checkbox
              v-model="f.options.animated"
            >
              {{ $t('kind.number.progress.animated') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          sm="3"
        >
          <b-form-group
            v-if="f.options.showValue"
            class="mb-0"
          >
            <b-form-checkbox
              v-model="f.options.showRelative"
              class="mb-2"
            >
              {{ $t('kind.number.progress.show.relative') }}
            </b-form-checkbox>

            <b-form-checkbox
              v-model="f.options.showProgress"
            >
              {{ $t('kind.number.progress.show.progress') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
        >
          <b-form-group>
            <template #label>
              <div
                class="d-flex align-items-center text-primary"
              >
                {{ $t('kind.number.progress.thresholds.label') }}
                <b-button
                  variant="link"
                  class="text-decoration-none ml-1"
                  @click="addThreshold()"
                >
                  {{ $t('general:label.add-with-plus') }}
                </b-button>
              </div>

              <small
                class="text-muted"
              >
                {{ $t('kind.number.progress.thresholds.description') }}
              </small>
            </template>

            <b-row
              v-for="(t, i) in field.options.thresholds"
              :key="i"
              align-v="center"
              :class="{ 'mt-2': i }"
            >
              <b-col>
                <b-input-group
                  append="%"
                >
                  <b-form-input
                    v-model="t.value"
                    :placeholder="'Threshold'"
                    type="number"
                    number
                  />
                </b-input-group>
              </b-col>

              <b-col
                class="d-flex align-items-center justify-content-center"
              >
                <b-form-select
                  v-model="t.variant"
                  :options="variants"
                />

                <font-awesome-icon
                  :icon="['fas', 'times']"
                  class="pointer text-danger ml-3"
                  @click="removeThreshold(i)"
                />
              </b-col>
            </b-row>
          </b-form-group>
        </b-col>
      </template>
    </b-row>

    <hr>

    <b-form-group
      :label=" $t('kind.number.liveExample')"
      label-class="text-primary"
      class="mb-0 w-100"
    >
      <b-row
        align-v="center"
      >
        <b-col
          cols="12"
          sm="6"
        >
          <b-form-input
            v-model="liveExample"
            type="number"
            number
          />
        </b-col>

        <b-col
          cols="12"
          sm="6"
        >
          <field-viewer
            value-only
            v-bind="mock"
          />
        </b-col>
      </b-row>
    </b-form-group>
  </div>
</template>

<script>
import base from './base'
import FieldViewer from '../Viewer'
import { compose, validator } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    FieldViewer,
  },

  extends: base,

  data () {
    return {
      liveExample: undefined,

      displayOptions: [
        { text: this.$t('kind.number.displayType.number'), value: 'number' },
        { text: this.$t('kind.number.displayType.progress'), value: 'progress' },
      ],

      variants: [
        { text: this.$t('kind.number.progress.variants.primary'), value: 'primary' },
        { text: this.$t('kind.number.progress.variants.secondary'), value: 'secondary' },
        { text: this.$t('kind.number.progress.variants.success'), value: 'success' },
        { text: this.$t('kind.number.progress.variants.warning'), value: 'warning' },
        { text: this.$t('kind.number.progress.variants.danger'), value: 'danger' },
        { text: this.$t('kind.number.progress.variants.info'), value: 'info' },
        { text: this.$t('kind.number.progress.variants.light'), value: 'light' },
        { text: this.$t('kind.number.progress.variants.dark'), value: 'dark' },
      ],

      mock: {
        namespace: undefined,
        module: undefined,
        field: undefined,
        record: undefined,
        errors: new validator.Validated(),
      },
    }
  },

  watch: {
    'field.options.display': {
      handler (display) {
        this.liveExample = display === 'number' ? 123.45679 : 33
      },
    },

    'field.options': {
      deep: true,
      handler (options) {
        if (this.mock.field) {
          this.mock.field.options = options
          this.mock.record.values.mockField = Number(this.liveExample).toFixed(this.field.options.precision)
        }
      },
    },

    liveExample: {
      handler (value) {
        if (this.mock.field) {
          value = Number(value).toFixed(this.field.options.precision)
          this.mock.record.values.mockField = value
        }
      },
    },
  },

  created () {
    this.mock.namespace = this.namespace
    this.mock.field = compose.ModuleFieldMaker(this.field)
    this.mock.field.isMulti = false
    this.mock.field.apply({ name: 'mockField' })
    this.mock.module = new compose.Module({ fields: [this.mock.field] }, this.namespace)
    this.mock.record = new compose.Record(this.mock.module, { })
    this.liveExample = this.field.options.display === 'number' ? 123.45679 : 33.45679
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    addThreshold () {
      this.field.options.thresholds.push({ value: 0, variant: 'success' })
    },

    removeThreshold (index) {
      if (index > -1) {
        this.field.options.thresholds.splice(index, 1)
      }
    },

    setDefaultValues () {
      this.liveExample = undefined
      this.displayOptions = []
      this.variants = []
      this.mock = {}
    },
  },
}
</script>
