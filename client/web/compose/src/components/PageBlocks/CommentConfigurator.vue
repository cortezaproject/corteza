<template>
  <b-tab :title="$t('comment.label')">
    <b-form-group>
      <label>{{ $t('general.module') }}</label>
      <b-form-select
        v-model="options.moduleID"
        :options="moduleOptions"
        text-field="name"
        value-field="moduleID"
        required
      />
    </b-form-group>
    <div v-if="selectedModule">
      <b-form-group>
        <label>{{ $t('field.selector.available') }}</label>
        <div class="d-flex">
          <div class="border fields w-100 p-2">
            <div
              v-for="field in allFields"
              :key="field.name"
              class="field"
            >
              <span v-if="field.label">{{ field.label }} ({{ field.name }})</span>
              <span v-else>{{ field.name }}</span>
              <span class="small float-right">
                <span v-if="field.isSystem">{{ $t('field.selector.systemField') }}</span>
                <span v-else>{{ field.kind }}</span>
              </span>
            </div>
          </div>
        </div>
      </b-form-group>

      <b-form-group
        horizontal
        :label-cols="3"
        breakpoint="md"
        :label="$t('recordList.record.prefilterLabel')"
      >
        <b-form-textarea
          v-model.trim="options.filter"
          :value="true"
          :placeholder="$t('recordList.record.prefilterPlaceholder')"
        />
        <b-form-text>
          <i18next
            path="recordList.record.prefilterFootnote"
            tag="label"
          >
            <code>${recordID}</code>
            <code>${ownerID}</code>
            <code>${userID}</code>
          </i18next>
        </b-form-text>
      </b-form-group>

      <b-form-group
        horizontal
        :label-cols="3"
        breakpoint="md"
        :label="$t('comment.titleField.label')"
      >
        <b-form-select v-model="options.titleField">
          <option value="">
            {{ $t('general.label.none') }}
          </option>
          <option
            v-for="(field, index) in selectedModuleFieldsByType('String')"
            :key="index"
            :value="field.name"
          >
            {{ field.label || field.name }} ({{ field.kind }})
          </option>
        </b-form-select>
        <b-form-text>{{ $t('comment.titleField.footnote') }}</b-form-text>
      </b-form-group>

      <b-form-group
        horizontal
        :label-cols="3"
        breakpoint="md"
        :label="$t('comment.contentField.label')"
      >
        <b-form-select v-model="options.contentField">
          <option value="">
            {{ $t('general.label.none') }}
          </option>
          <option
            v-for="(field, index) in selectedModuleFieldsByType('String')"
            :key="index"
            :value="field.name"
          >
            {{ field.label || field.name }} ({{ field.kind }})
          </option>
        </b-form-select>
        <b-form-text class="text-secondary small">
          {{ $t('comment.contentField.footnote') }}
        </b-form-text>
      </b-form-group>
      <b-form-group
        horizontal
        :label-cols="3"
        breakpoint="md"
        :label="$t('comment.referenceField.label')"
      >
        <b-form-select v-model="options.referenceField">
          <option value="">
            {{ $t('general.label.none') }}
          </option>
          <option
            v-for="(field, index) in selectedModuleFieldsByType('Record')"
            :key="index"
            :value="field.name"
          >
            {{ field.label || field.name }} ({{ field.kind }})
          </option>
        </b-form-select>
        <b-form-text class="text-secondary small">
          {{ $t('comment.referenceField.footnote') }}
        </b-form-text>
      </b-form-group>
    </div>
    <b-form-group
      horizontal
      :label-cols="3"
      breakpoint="md"
      :label="$t('comment.sortDirection.label')"
    >
      <b-form-select v-model="options.sortDirection">
        <option
          v-for="(item, index) in sortDirections"
          :key="index"
          :value="item.value"
        >
          {{ item.label }}
        </option>
      </b-form-select>
      <b-form-text class="text-secondary small">
        {{ $t('comment.sortDirection.footnote') }}
      </b-form-text>
    </b-form-group>
  </b-tab>
</template>
<script>
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'CommentConfigurator',

  extends: base,
  data () {
    return {
      referenceList: [{ label: 'Record ID (recordID)', value: 'recordID' }, { label: 'Page ID (pageID)', value: 'pageID' }],
      sortDirections: [{ label: this.$t('comment.sortDirection.asc'), value: 'asc' }, { label: this.$t('comment.sortDirection.desc'), value: 'desc' }],
    }
  },
  computed: {
    ...mapGetters({
      modules: 'module/set',
    }),

    moduleOptions () {
      return [
        { moduleID: NoID, name: this.$t('general.label.none') },
        ...this.filterModulesByRecord,
      ]
    },

    filterModulesByRecord () {
      if (this.record) {
        return this.modules.filter(module => {
          return module.fields.some(f => {
            if (f.kind === 'Record') {
              if (f.options.moduleID === this.options.moduleID) {
                return false
              }
            }
            return true
          })
        })
      }
      return this.modules
    },

    selectedModule () {
      return this.modules.find(m => m.moduleID === this.options.moduleID)
    },

    selectedModuleFields () {
      if (this.selectedModule) {
        return [...this.selectedModule.fields].sort((a, b) => a.label.localeCompare(b.label))
      }
      return []
    },

    allFields () {
      if (this.options.moduleID) {
        return [
          ...this.selectedModuleFields,
          ...this.selectedModule.systemFields().map(sf => {
            sf.label = this.$t(`field:system.${sf.name}`)
            return sf
          }),
        ]
      }
      return []
    },
  },

  watch: {
    'options.moduleID': {
      handler () {
        this.options.titleField = ''
        this.options.contentField = ''
        this.options.referenceField = ''
        this.selectedModuleFields.forEach(f => {
          if (f.name === 'Reference') {
            this.options.referenceField = 'Reference'
          }
          if (f.name === 'Content') {
            this.options.contentField = 'Content'
          }
        })
      },
    },
  },

  created () {
    if (!this.options.sortDirection) {
      this.options.sortDirection = 'desc'
    }
  },
  methods: {
    selectedModuleFieldsByType (type) {
      return (this.selectedModuleFields || []).filter((f) => {
        return f.kind === type
      })
    },
  },
}
</script>
<style lang="scss" scoped>
.fields {
  height: 150px;
  overflow-y: auto;
  cursor: default;
}
</style>
