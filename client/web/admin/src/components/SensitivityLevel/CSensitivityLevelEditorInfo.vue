<template>
  <b-card
    data-test-id="card-sens-lvl-info"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <b-form
      @submit.prevent="$emit('submit', sensitivityLevel)"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('name')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="sensitivityLevel.meta.name"
              data-test-id="input-name"
              required
              :state="nameState"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('handle.label')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="sensitivityLevel.handle"
              :placeholder="$t('handle.placeholder')"
              :state="handleState"
            />

            <b-form-invalid-feedback :state="handleState">
              {{ $t('handle.invalid-characters') }}
            </b-form-invalid-feedback>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('level', sensitivityLevel)"
          >
            <b-form-input
              v-model="sensitivityLevel.level"
              number
              type="range"
              min="1"
              max="10"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('description')"
          >
            <b-form-textarea
              v-model="sensitivityLevel.meta.description"
            />
          </b-form-group>
        </b-col>
      </b-row>
      <c-system-fields
        :resource="sensitivityLevel"
      />
      <!--
        include hidden input to enable
        trigger submit event w/ ENTER
      -->
      <input
        type="submit"
        class="d-none"
        :disabled="saveDisabled"
      >
    </b-form>

    <template #footer>
      <c-input-confirm
        v-if="sensitivityLevel && sensitivityLevel.sensitivityLevelID"
        variant="danger"
        size="md"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </c-input-confirm>

      <c-button-submit
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', sensitivityLevel)"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'

export default {
  name: 'CSensitivityLevelEditorInfo',

  i18nOptions: {
    namespaces: 'system.sensitivityLevel',
    keyPrefix: 'editor.info',
  },

  props: {
    sensitivityLevel: {
      type: Object,
      required: true,
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
    fresh () {
      return !this.sensitivityLevel.sensitivityLevelID || this.sensitivityLevel.sensitivityLevelID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : true
    },

    nameState () {
      return this.sensitivityLevel.meta.name ? null : false
    },

    handleState () {
      return handle.handleState(this.sensitivityLevel.handle)
    },

    saveDisabled () {
      return !this.editable || [this.nameState, this.handleState].includes(false)
    },

    getDeleteStatus () {
      return this.sensitivityLevel.deletedAt ? this.$t('undelete') : this.$t('delete')
    },
  },
}
</script>
