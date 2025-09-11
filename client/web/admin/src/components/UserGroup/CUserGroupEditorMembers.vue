<template>
  <b-card
    data-test-id="card-user-group-edit-members"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <b-form @submit.prevent="$emit('submit')">
      <c-member-picker
        :value="value"
        no-remove
        @input="$emit('input', $event)"
      />
    </b-form>

    <template #header>
      <h4 class="m-0">
        {{ $t("title") }}
      </h4>
    </template>

    <template #footer>
      <c-button-submit
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit')"
      />
    </template>
  </b-card>
</template>

<script>
import CMemberPicker from 'corteza-webapp-admin/src/components/CMemberPicker'

export default {
  i18nOptions: {
    namespaces: 'system.user-groups',
    keyPrefix: 'editor.members',
  },

  components: {
    CMemberPicker,
  },

  props: {
    value: {
      type: Array,
      required: true,
      default: () => [],
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },
  },
}
</script>
