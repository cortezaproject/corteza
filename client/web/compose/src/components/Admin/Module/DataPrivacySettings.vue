<template>
  <div
    v-if="resource && connection"
  >
    <b-form-group
      :label="translations.sensitivity.label"
      :description="translations.sensitivity.description"
      label-class="text-primary"
    >
      <c-sensitivity-level-picker
        v-model="resource.config.privacy.sensitivityLevelID"
        :options="sensitivityLevels"
        :placeholder="translations.sensitivity.placeholder"
        :max-level="maxLevel"
        :disabled="processing"
      />
    </b-form-group>
    <b-form-group
      :label="translations.usage.label"
      label-class="text-primary"
    >
      <b-textarea
        v-model="resource.config.privacy.usageDisclosure"
      />
    </b-form-group>
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CSensitivityLevelPicker } = components

export default {
  components: {
    CSensitivityLevelPicker,
  },

  props: {
    resource: {
      type: Object,
      required: true,
    },

    connection: {
      type: Object,
      required: true,
    },

    // ID of sensitivityLevel with the maximum allowed level
    maxLevel: {
      type: String,
      default: undefined,
    },

    sensitivityLevels: {
      type: Array,
      required: true,
    },

    translations: {
      type: Object,
      required: true,
    },
  },

  data () {
    return {
      processing: false,
    }
  },
}
</script>
