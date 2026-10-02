<template>
  <b-card
    :title="sourceLabel"
    footer-class="text-right pt-0"
    class="border"
  >
    <p
      v-if="workflow && workflow.meta"
      class="mb-3"
    >
      {{ workflow.meta.description || $t('noDescription') }}

      <var
        v-if="trigger"
      >
        {{ $t('stepID', { stepID: trigger.stepID }) }}
      </var>
    </p>

    <p
      v-else-if="script"
      class="mb-3"
    >
      {{ script.description || $t('noDescription') }}
    </p>

    <b-form-group
      :label="$t('buttonLabel')"
      label-class="text-primary"
    >
      <c-input-expression
        v-model="button.label"
        auto-complete
        :suggestion-params="recordAutoCompleteParams"
        :page="page"
      />
      <i18next
        path="block:interpolationFootnote"
        tag="small"
        class="text-muted"
      >
        <code>${record.values.fieldName}</code>
        <code>${recordID}</code>
        <code>${ownerID}</code>
        <span><code>${userID}</code>, <code>${user.name}</code></span>
      </i18next>
    </b-form-group>

    <b-form-group
      :label="$t('buttonVariant')"
      label-class="text-primary"
    >
      <b-select
        v-model="button.variant"
        class="w-100"
      >
        <b-select-option
          v-for="({ variant, label }) in variants"
          :key="variant"
          :value="variant"
        >
          {{ label }}
        </b-select-option>
      </b-select>
    </b-form-group>

    <b-form-group
      label-class="text-primary"
    >
      <template #label>
        <div class="d-flex align-items-center">
          {{ $t('buttonVisibility.label') }}

          <c-hint
            :tooltip="$t('buttonVisibility.tooltip')"
            icon-class="text-warning"
          />
        </div>
      </template>

      <b-input-group>
        <b-input-group-prepend>
          <b-input-group-text variant="extra-light">
            ƒ
          </b-input-group-text>
        </b-input-group-prepend>

        <c-input-expression
          v-model="visibilityExpression"
          auto-complete
          :placeholder="$t('buttonVisibility.placeholder')"
          :suggestion-params="visibilityAutoCompleteParams"
          class="flex-grow-1"
        />
      </b-input-group>

      <i18next
        v-if="isRecordPage"
        path="block:general.visibility.condition.description.record-page"
        tag="small"
        class="text-muted"
      >
        <code>record.values.fieldName</code>
        <code>user.(userID/email...)</code>
        <code>screen.(width/height/userAgent/breakpoint)</code>
        <code>user.userID == record.createdBy</code>
        <code>record.values.fieldName == "value"</code>
        <code>screen.width &lt; 1024</code>
      </i18next>

      <i18next
        v-else
        path="block:general.visibility.condition.description.non-record-page"
        tag="small"
        class="text-muted"
      >
        <code>user.(userID/email...)</code>
        <code>screen.(width/height/userAgent/breakpoint)</code>
        <code>user.userID == "123456789"</code>
        <code>screen.width &lt; 1024</code>
      </i18next>
    </b-form-group>

    <template #footer>
      <c-input-confirm
        show-icon
        @confirmed="$emit('delete', button)"
      />
    </template>
  </b-card>
</template>
<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
import autocomplete from 'corteza-webapp-compose/src/mixins/autocomplete.js'

const { CInputExpression } = components

export default {
  i18nOptions: {
    namespaces: 'block',
    keyPrefix: 'automation',
  },

  components: {
    CInputExpression,
  },

  mixins: [autocomplete],

  props: {
    button: {
      type: Object,
      required: true,
    },

    script: {
      type: Object,
      required: false,
      default: undefined,
    },

    trigger: {
      type: Object,
      required: false,
      default: undefined,
    },

    page: {
      type: compose.Page,
      required: true,
    },

    block: {
      type: compose.PageBlock,
      required: true,
    },

    module: {
      type: compose.Module,
      required: false,
      default: undefined,
    },

    record: {
      type: compose.Record,
      required: false,
      default: undefined,
    },
  },

  computed: {
    visibilityExpression: {
      get () {
        return (this.button.visibility || {}).expression || ''
      },

      set (expression) {
        this.$set(this.button, 'visibility', { ...(this.button.visibility || {}), expression })
      },
    },

    variants () {
      return [
        'primary',
        'secondary',
        'light',
        'dark',
        'success',
        'danger',
        'warning',
      ].map(variant => ({ variant, label: this.$t(`variants.${variant}`) }))
    },

    sourceLabel () {
      if (this.workflow) {
        return this.workflow.meta.name || this.$t('noLabel')
      } else if (this.button.script) {
        return this.button.script
      }

      return this.$t('dummyButtonLabel')
    },

    workflow () {
      return this.trigger ? this.trigger.workflow : undefined
    },

    isNew () {
      return this.block.blockID === NoID
    },
  },
}
</script>
