<template>
  <b-card
    class="shadow-sm"
    :title="$t('title')"
  >
    <b-row>
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('form.name.label')"
          :description="$t('form.name.description')"
          class="mb-3 text-primary"
        >
          <b-form-input
            v-model="connection.meta.name"
            :disabled="disabled"
            :placeholder="$t('form.name.placeholder')"
          />
        </b-form-group>
      </b-col>
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('form.handle.label')"
          :description="$t('form.handle.description')"
          class="mb-3 text-primary"
        >
          <b-form-input
            v-model="connection.handle"
            :disabled="isPrimary || disabled"
            :placeholder="$t('form.handle.placeholder')"
            :state="handleState"
          />
          <b-form-invalid-feedback :state="handleState">
            {{ $t('form.handle.invalid-characters') }}
          </b-form-invalid-feedback>
        </b-form-group>
      </b-col>
    </b-row>
    <b-row>
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('form.location-name.label')"
          :description="$t('form.location-name.description')"
        >
          <b-form-input
            v-model="connection.meta.location.properties.name"
            :disabled="disabled"
            :placeholder="$t('form.location-name.placeholder')"
          />
        </b-form-group>
      </b-col>
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :description="$t('form.location-geometry.description')"
        >
          <label
            class="d-flex align-items-center text-primary"
          >
            {{ $t('form.location-geometry.label') }}
            <c-location
              v-if="!disabled"
              :value="connection.meta.location.geometry.coordinates || []"
              :placeholder="$t('form.location-geometry.placeholder')"
              class="ml-1"
              editable
              @input="connection.meta.location.geometry.coordinates = $event"
            />
          </label>

          <div
            class="mt-2 d-flex align-items-center"
          >
            <code
              v-if="locationCoordinates"
            >
              {{ locationCoordinates }}
            </code>

            <span v-else>
              -
            </span>
          </div>
        </b-form-group>
      </b-col>
    </b-row>

    <b-row>
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('form.ownership.label')"
          :description="$t('form.ownership.description')"
          class="mb-3 text-primary"
        >
          <b-form-input
            v-model="connection.meta.ownership"
            :disabled="disabled"
            :placeholder="$t('form.ownership.placeholder')"
          />
        </b-form-group>
      </b-col>
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('form.sensitivity-level.label')"
          :description="$t('form.sensitivity-level.description')"
          class="text-primary"
        >
          <c-sensitivity-level-picker
            v-model="connection.config.privacy.sensitivityLevelID"
            :options="sensitivityLevels"
            :disabled="disabled"
            :placeholder="$t('form.sensitivity-level.placeholder')"
          />
        </b-form-group>
      </b-col>
    </b-row>
  </b-card>
</template>

<script>
import { components, handle } from '@cortezaproject/corteza-vue'
import CLocation from 'corteza-webapp-admin/src/components/CLocation'
const { CSensitivityLevelPicker } = components

export default {
  i18nOptions: {
    namespaces: 'system.connections',
    keyPrefix: 'editor.basic',
  },

  components: {
    CLocation,
    CSensitivityLevelPicker,
  },

  props: {
    disabled: { type: Boolean, default: false },
    isPrimary: { type: Boolean, required: true },

    connection: {
      type: Object,
      required: true,
    },

    sensitivityLevels: {
      type: Array,
      required: true,
    },
  },

  computed: {
    handleState () {
      return handle.handleState(this.connection.handle)
    },

    locationCoordinates () {
      const { coordinates: cc } = this.connection.meta.location.geometry

      if (cc && Array.isArray(cc) && cc.length === 2) {
        return cc.map(c => c.toFixed(7)).join(', ')
      }

      return ''
    },
  },
}
</script>
