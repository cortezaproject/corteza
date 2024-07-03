<template>
  <b-card
    data-test-id="card-role-info"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <b-form
      @submit.prevent="submit()"
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
              v-model="role.name"
              data-test-id="input-name"
              :state="nameState"
              :disabled="!editable"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('handle')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="role.handle"
              data-test-id="input-handle"
              :state="handleState"
              :disabled="!editable"
              :placeholder="$t('placeholder-handle')"
            />

            <b-form-invalid-feedback :state="handleState">
              {{ $t('invalid-handle-characters') }}
            </b-form-invalid-feedback>
          </b-form-group>
        </b-col>

        <b-col cols="12">
          <b-form-group
            v-if="role.meta"
            :label="$t('description')"
            label-class="text-primary"
          >
            <b-form-textarea
              v-model="role.meta.description"
              data-test-id="textarea-description"
              :disabled="!editable"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('context.label')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="isContextual"
              data-test-id="checkbox-is-contextual"
              switch
              :labels="checkboxLabel"
              :disabled="!editable"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <b-row
        v-if="isContextual"
        class="my-3"
      >
        <b-col
          cols="12"
        >
          <b-form-group
            :label="$t('context.expression-label')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="role.meta.context.expr"
              data-test-id="input-expression"
              :disabled="!editable"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
        >
          <b-form-group
            :label="$t('context.resource-types-label')"
            label-class="text-primary"
          >
            <b-checkbox
              v-for="(resourceType, i) in resourceTypeOptions"
              :key="i"
              v-model="role.meta.context.resourceTypes"
              :data-test-id="`checkbox-resource-type-${resourceType.text}`"
              :value="resourceType.value"
              :disabled="!editable"
            >
              {{ resourceType.text }}
            </b-checkbox>
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :resource="role"
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

    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <template #footer>
      <c-input-confirm
        v-if="!fresh && editable && role.canDeleteRole && !isDataPrivacyOfficer"
        :data-test-id="deletedButtonStatusCypressId"
        variant="danger"
        size="md"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </c-input-confirm>

      <c-input-confirm
        v-if="!fresh && editable && !isDataPrivacyOfficer"
        :data-test-id="archivedButtonStatusCypressId"
        variant="secondary"
        size="md"
        @confirmed="$emit('status')"
      >
        {{ getArchiveStatus }}
      </c-input-confirm>

      <c-button-submit
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="submit()"
      />
    </template>
  </b-card>
</template>

<script>
import { system, NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'

export default {
  name: 'CRoleEditorInfo',

  i18nOptions: {
    namespaces: 'system.roles',
    keyPrefix: 'editor.info',
  },

  props: {
    role: {
      type: system.Role,
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

    isContext: {
      type: Boolean,
      required: true,
    },

    canCreate: {
      type: Boolean,
      required: true,
    },
  },

  data () {
    return {
      checkboxLabel: {
        on: this.$t('general:label.general.yes'),
        off: this.$t('general:label.general.no'),
      },
    }
  },

  computed: {
    isContextual: {
      get () {
        return this.isContext
      },

      set (isContext) {
        this.$emit('update:is-context', isContext)
      },
    },

    saveDisabled () {
      return !this.editable || [this.nameState, this.handleState].includes(false)
    },

    /**
     * Returns true if role is not saved yet (without ID)
     * @returns {boolean}
     */
    fresh () {
      return !this.role.roleID || this.role.roleID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : !this.role.isSystem && this.role.canUpdateRole
    },

    // At least 1 character
    nameState () {
      const { name } = this.role
      return name ? null : false
    },

    // 2+ alpha-numeric + _
    handleState () {
      return handle.handleState(this.role.handle)
    },

    getDeleteStatus () {
      return this.role.deletedAt ? this.$t('undelete') : this.$t('delete')
    },

    getArchiveStatus () {
      return this.role.archivedAt ? this.$t('unarchive') : this.$t('archive')
    },

    resourceTypeOptions () {
      return this.resourceTypes.map(value => ({
        // @todo use translation facility to generate resource type option labels
        text: value.replace('corteza::', ''),
        value,
      }))
    },

    resourceTypes () {
      // @todo this should be fetched from the backend
      return [
        // 'corteza::system:application',
        'corteza::system:auth-client',
        'corteza::system:role',
        // 'corteza::system:template',
        'corteza::system:user',
        // 'corteza::compose:chart',
        // 'corteza::compose:module-field',
        'corteza::compose:module',
        'corteza::compose:namespace',
        'corteza::compose:page',
        'corteza::compose:record',
        'corteza::automation:workflow',
        // 'corteza::federation:exposed-module',
        // 'corteza::federation:node',
        // 'corteza::federation:shared-module',
      ]
    },

    deletedButtonStatusCypressId () {
      return `button-${this.getDeleteStatus.toLowerCase()}`
    },

    archivedButtonStatusCypressId () {
      return `button-${this.getArchiveStatus.toLowerCase()}`
    },

    isDataPrivacyOfficer () {
      return this.role.handle === 'data-privacy-officer'
    },
  },

  methods: {
    submit () {
      if (!this.isContextual && this.role.isContext) {
        // if checkbox was unchecked on submit, purge meta before submit
        this.role.meta.context.resourceTypes = []
        this.role.meta.context.expr = ''
      }
      this.$emit('submit', this.role)
    },
  },
}
</script>
