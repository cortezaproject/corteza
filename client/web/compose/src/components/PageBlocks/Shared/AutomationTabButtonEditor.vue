<template>
  <b-card
    :header="workflow ? $t('editTitle.workflow') : $t('editTitle.script')"
    footer-class="text-right"
    class="border"
  >
    <b-form-group
      :label="$t('buttonLabel')"
      label-class="text-primary"
    >
      <c-input-expression
        v-model="button.label"
        auto-complete
        lang="javascript"
        :suggestion-params="recordAutoCompleteParams"
        :page="page"
        height="2.375rem"
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
        size="sm"
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

    <div
      v-if="workflow"
    >
      <h5>
        {{ workflow.meta.name || $t('noLabel') }}
      </h5>
    </div>

    <code
      v-else-if="button.script"
    >
      {{ button.script }}
    </code>

    <b-alert
      v-else
      show
      variant="warning"
    >
      {{ $t('noScript' ) }}
    </b-alert>

    <p
      v-if="workflow && workflow.meta"
      class="my-2"
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
      class="my-2"
    >
      {{ script.description || $t('noDescription') }}
    </p>

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

  extends: autocomplete,

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
  },

  computed: {
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

    workflow () {
      return this.trigger ? this.trigger.workflow : undefined
    },

    isNew () {
      return this.block.blockID === NoID
    },
  },
}
</script>
