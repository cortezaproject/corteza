<template>
  <b-tab
    :title="$t('progress.tab.label')"
    class="vh-100"
  >
    <template>
      <h5 class="text-primary">
        {{ $t('progress.value.label') }}
      </h5>

      <b-row>
        <b-col
          cols="12"
        >
          <b-form-group
            :label="$t('progress.module.label')"
          >
            <vue-select
              v-model="options.value.moduleID"
              label="name"
              :placeholder="$t('progress.module.select')"
              :options="modules"
              :reduce="m => m.moduleID"
              class="bg-white"
            />
          </b-form-group>
        </b-col>

        <template v-if="options.value.moduleID">
          <b-col
            cols="12"
          >
            <b-form-group :label="$t('metric.edit.filterLabel')">
              <b-form-textarea
                v-model="options.value.filter"
                placeholder="(A > B) OR (A < C)"
                class="mb-1"
              />
              <b-form-text>
                <i18next
                  path="metric.edit.filterFootnote"
                  tag="label"
                >
                  <code>${recordID}</code>
                  <code>${ownerID}</code>
                  <code>${userID}</code>
                </i18next>
              </b-form-text>
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            sm="6"
          >
            <b-form-group
              :label="$t('progress.field.label')"
            >
              <vue-select
                v-model="options.value.field"
                :placeholder="$t('progress.field.select')"
                :options="valueModuleFields"
                :reduce="f => f.name"
                class="bg-white"
                @input="fieldChanged($event, options.value)"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            sm="6"
          >
            <b-form-group
              :label="$t('progress.aggregate.label')"
            >
              <vue-select
                v-model="options.value.operation"
                label="name"
                :disabled="!options.value.field || options.value.field === 'count'"
                :placeholder="$t('progress.aggregate.select')"
                :options="aggregationOperations"
                :reduce="a => a.operation"
                class="bg-white"
              />
            </b-form-group>
          </b-col>
        </template>
      </b-row>
    </template>

    <hr>

    <template>
      <h5 class="text-primary">
        {{ $t('progress.value.max') }}
      </h5>

      <b-row>
        <b-col
          cols="12"
        >
          <b-form-group
            :label="$t('progress.module.label')"
          >
            <vue-select
              v-model="options.maxValue.moduleID"
              label="name"
              :placeholder="$t('progress.module.select')"
              :options="modules"
              :reduce="m => m.moduleID"
              class="bg-white"
            />
          </b-form-group>
        </b-col>

        <template v-if="options.maxValue.moduleID">
          <b-col
            cols="12"
          >
            <b-form-group :label="$t('metric.edit.filterLabel')">
              <b-form-textarea
                v-model="options.maxValue.filter"
                placeholder="(A > B) OR (A < C)"
                class="mb-1"
              />
              <b-form-text>
                <i18next
                  path="metric.edit.filterFootnote"
                  tag="label"
                >
                  <code>${recordID}</code>
                  <code>${ownerID}</code>
                  <code>${userID}</code>
                </i18next>
              </b-form-text>
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            sm="6"
          >
            <b-form-group
              :label="$t('progress.field.label')"
            >
              <vue-select
                v-model="options.maxValue.field"
                :placeholder="$t('progress.field.select')"
                :options="maxValueModuleFields"
                :reduce="f => f.name"
                class="bg-white"
                @input="fieldChanged($event, options.maxValue)"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            sm="6"
          >
            <b-form-group
              :label="$t('progress.aggregate.label')"
            >
              <vue-select
                v-model="options.maxValue.operation"
                label="name"
                :disabled="!options.maxValue.field || options.maxValue.field === 'count'"
                :placeholder="$t('progress.aggregate.select')"
                :options="aggregationOperations"
                :reduce="a => a.operation"
                class="bg-white"
              />
            </b-form-group>
          </b-col>
        </template>
      </b-row>
    </template>

    <hr>

    <template>
      <h5 class="text-primary">
        {{ $t('progress.display-options') }}
      </h5>

      <b-row
        align-v="center"
      >
        <b-col
          cols="12"
          sm="6"
        >
          <b-form-group
            :label="$t('progress.default-variant')"
          >
            <b-form-select
              v-model="options.display.variant"
              :options="variants"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="6"
          sm="3"
          class="mb-2 mb-sm-0"
        >
          <b-form-group
            class="mb-0"
          >
            <b-form-checkbox
              v-model="options.display.showValue"
              class="mb-2"
            >
              {{ $t('progress.show.value') }}
            </b-form-checkbox>
            <b-form-checkbox
              v-model="options.display.animated"
            >
              {{ $t('progress.animated') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>

        <b-col
          cols="6"
          sm="3"
          class="mb-2 mb-sm-0"
        >
          <b-form-group
            v-if="options.display.showValue"
            class="mb-0"
          >
            <b-form-checkbox
              v-model="options.display.showRelative"
              class="mb-2"
            >
              {{ $t('progress.show.relative') }}
            </b-form-checkbox>
            <b-form-checkbox
              v-model="options.display.showProgress"
            >
              {{ $t('progress.show.progress') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
        >
          <b-form-group>
            <template #label>
              <div
                class="d-flex align-items-center"
              >
                {{ $t('progress.thresholds') }}
                <b-button
                  variant="link"
                  class="text-decoration-none ml-1"
                  @click="addThreshold()"
                >
                  {{ $t('progress.add') }}
                </b-button>
              </div>

              <small
                class="text-muted"
              >
                {{ $t('progress.threshold.variant') }}
              </small>
            </template>
            <b-row
              v-for="(t, i) in options.display.thresholds"
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
                    :placeholder="$t('progress.threshold.label')"
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
      </b-row>

      <!-- Don't remove, makes sure preview doesn't cover content -->
      <div class="w-501 my-5 py-5 py-sm-3" />

      <template>
        <div
          class="preview bg-white position-absolute p-3"
        >
          <h6 class="text-primary">
            {{ $t('progress.preview') }}
          </h6>

          <b-progress
            :max="preview.max"
            height="1.5rem"
            class="mb-2 bg-light"
          >
            <b-progress-bar
              :value="preview.value"
              :striped="options.display.striped"
              :animated="options.display.animated"
              :variant="progressVariant"
            >
              {{ progressLabel }}
            </b-progress-bar>
          </b-progress>

          <b-row>
            <b-col
              cols="12"
              sm="6"
            >
              {{ $t('progress.value.label') }}
              <b-form-input
                v-model="preview.value"
                :placeholder="$t('progress.value.label')"
                type="number"
                size="sm"
              />
            </b-col>

            <b-col
              cols="12"
              sm="6"
            >
              {{ $t('progress.value.max') }}
              <b-form-input
                v-model="preview.max"
                :placeholder="$t('progress.value.max')"
                type="number"
                size="sm"
              />
            </b-col>
          </b-row>
        </div>
      </template>
    </template>
  </b-tab>
</template>

<script>
import base from './base'
import { mapGetters } from 'vuex'
import { VueSelect } from 'vue-select'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'ProgressConfigurator',

  components: {
    VueSelect,
  },

  extends: base,

  data () {
    return {
      preview: {
        value: 15,
        max: 50,
      },

      aggregationOperations: [
        {
          name: this.$t('metric.edit.operationSum'),
          operation: 'sum',
        },
        {
          name: this.$t('metric.edit.operationMax'),
          operation: 'max',
        },
        {
          name: this.$t('metric.edit.operationMin'),
          operation: 'min',
        },
        {
          name: this.$t('metric.edit.operationAvg'),
          operation: 'avg',
        },
      ],

      variants: [
        { text: this.$t('progress.variant.primary'), value: 'primary' },
        { text: this.$t('progress.variant.secondary'), value: 'secondary' },
        { text: this.$t('progress.variant.success'), value: 'success' },
        { text: this.$t('progress.variant.warning'), value: 'warning' },
        { text: this.$t('progress.variant.danger'), value: 'danger' },
        { text: this.$t('progress.variant.info'), value: 'info' },
        { text: this.$t('progress.variant.light'), value: 'light' },
        { text: this.$t('progress.variant.dark'), value: 'dark' },
      ],
    }
  },

  computed: {
    ...mapGetters({
      modules: 'module/set',
      moduleByID: 'module/getByID',
    }),

    sharedModuleFields () {
      return [
        { name: 'count', label: 'Count' },
      ]
    },

    valueModuleFields () {
      return [
        ...this.sharedModuleFields,
        ...this.moduleByID(this.options.value.moduleID).fields.filter(f => f.kind === 'Number').sort((a, b) => a.name.localeCompare(b.name)),
      ]
    },

    maxValueModuleFields () {
      return [
        ...this.sharedModuleFields,
        ...this.moduleByID(this.options.maxValue.moduleID).fields.filter(f => f.kind === 'Number').sort((a, b) => a.name.localeCompare(b.name)),
      ]
    },

    progress () {
      const { value = 0, max = 100 } = this.preview
      return 100 * (value / max)
    },

    progressLabel () {
      let { value, max } = this.preview
      const { showValue, showRelative, showProgress } = this.options.display || {}

      if (!showValue) {
        return
      }

      if (showRelative) {
        // https://stackoverflow.com/a/21907972/17926309
        value = `${Math.round(((value / max) * 100) * 100) / 100}%`
      }

      if (showProgress) {
        value = `${value} / ${showRelative ? '100' : max}${showRelative ? '%' : ''}`
      }

      return value
    },

    progressVariant () {
      const { variant } = this.options.display || {}
      let progressVariant = variant

      if (this.options.display.thresholds.length) {
        const sortedVariants = [...this.options.display.thresholds].filter(t => t.value >= 0).sort((a, b) => b.value - a.value)
        const { variant } = sortedVariants.find(t => this.progress >= t.value) || {}
        progressVariant = variant || progressVariant
      }

      return progressVariant
    },
  },

  methods: {
    addThreshold () {
      this.options.display.thresholds.push({ value: 0, variant: 'success' })
    },

    removeThreshold (index) {
      if (index > -1) {
        this.options.display.thresholds.splice(index, 1)
      }
    },

    fieldChanged (value, optionsType) {
      if (!value) {
        optionsType.operation = ''
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.preview {
  bottom: 0;
  left: 0;
  z-index: 2;
  width: 100%;
  box-shadow: 0 -0.25rem 1rem rgb(0 0 0 / 15%);
}
</style>
