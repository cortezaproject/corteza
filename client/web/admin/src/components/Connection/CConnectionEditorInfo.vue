<template>
  <b-card
    data-test-id="card-connection-settings"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <b-row>
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('form.name.label')"
          :description="$t('form.name.description')"
          label-class="text-primary"
          class="mb-3"
        >
          <b-form-input
            v-model="connection.meta.name"
            required
            :placeholder="$t('form.name.placeholder')"
            :state="nameState"
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
          label-class="text-primary"
          class="mb-3"
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
          label-class="text-primary"
        >
          <b-form-input
            v-model="connection.meta.location.properties.name"
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
          label-class="text-primary"
          class="mb-3"
        >
          <b-form-input
            v-model="connection.meta.ownership"
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
          label-class="text-primary"
        >
          <c-sensitivity-level-picker
            v-model="connection.config.privacy.sensitivityLevelID"
            :options="sensitivityLevels"
            :placeholder="$t('form.sensitivity-level.placeholder')"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <template #footer>
      <c-input-confirm
        v-if="!fresh && !isPrimary && !disabled"
        variant="danger"
        size="md"
        @confirmed="$emit('delete')"
      >
        {{ connection.deletedAt ? $t('general:label.undelete') : $t('general:label.delete') }}
      </c-input-confirm>

      <c-button-submit
        :disabled="disabled || saveDisabled"
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
import { NoID } from '@cortezaproject/corteza-js'
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
    connection: {
      type: Object,
      required: true,
    },

    sensitivityLevels: {
      type: Array,
      required: true,
    },

    disabled: {
      type: Boolean,
      default: false,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    canCreate: {
      type: Boolean,
      required: true,
    },
  },

  computed: {
    isPrimary () {
      return this.connection.type === 'corteza::system:primary-dal-connection'
    },

    fresh () {
      return !this.connection.connectionID || this.connection.connectionID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : true
    },

    nameState () {
      return this.connection.meta.name ? null : false
    },

    handleState () {
      return handle.handleState(this.connection.handle)
    },

    saveDisabled () {
      return !this.editable || [this.nameState, this.handleState].includes(false)
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
