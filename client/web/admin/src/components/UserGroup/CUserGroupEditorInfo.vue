<template>
  <b-card
    data-test-id="card-user-group-info"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <b-form
      @submit.prevent="$emit('submit', userGroup)"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('meta.short')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="userGroup.meta.short"
              data-test-id="input-name"
              required
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('handle')"
            :class="{ 'mb-0': !userGroup.userGroupID }"
            label-class="text-primary"
          >
            <b-form-input
              v-model="userGroup.handle"
              data-test-id="input-handle"
              :placeholder="$t('placeholder-handle')"
              :state="handleState"
            />
            <b-form-invalid-feedback :state="handleState">
              {{ $t('invalid-handle-characters') }}
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
            :label="$t('parent')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="userGroup.selfID"
              data-test-id="select-parent"
              :options="parentUserGroups"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <b-row>
        <b-col
          cols="12"
        >
          <b-form-group
            :label="$t('meta.description')"
            label-class="text-primary"
          >
            <b-form-textarea
              v-model="userGroup.meta.description"
              data-test-id="textarea-description"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :id="userGroup.userGroupID"
        :resource="userGroup"
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
        v-if="!fresh && userGroup.canDeleteUserGroup"
        :data-test-id="deletedButtonStatusCypressId"
        :text="getDeleteStatus"
        variant="danger"
        size="md"
        @confirmed="$emit('delete')"
      />

      <c-corredor-manual-buttons
        ui-page="user-group/editor"
        ui-slot="infoFooter"
        resource-type="system:user-group"
        default-variant="light"
        @click="dispatchCortezaSystemUserGroupEvent($event, { userGroup })"
      />

      <c-button-submit
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', userGroup)"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'
import { getSystemFields } from 'corteza-webapp-admin/src/lib/sysFields'

export default {
  name: 'CUserEditorInfo',

  i18nOptions: {
    namespaces: 'system.user-groups',
    keyPrefix: 'editor.info',
  },

  props: {
    userGroup: {
      type: Object,
      required: true,
    },

    parentUserGroups: {
      type: Array,
      required: false,
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

    canCreate: {
      type: Boolean,
      required: true,
    },
  },

  computed: {
    getDeleteStatus () {
      return this.userGroup.deletedAt ? this.$t('undelete') : this.$t('delete')
    },

    userGroupID () {
      if (this.$auth.userGroup) {
        return this.$auth.userGroup.userGroupID
      }
      return undefined
    },

    fresh () {
      return !this.userGroup.userGroupID || this.userGroup.userGroupID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : this.userGroup.canUpdateUserGroup
    },

    handleState () {
      return handle.handleState(this.userGroup.handle)
    },

    saveDisabled () {
      return !this.editable || [this.handleState].includes(false)
    },

    deletedButtonStatusCypressId () {
      return `button-${this.getDeleteStatus.toLowerCase()}`
    },

    suspendButtonStatusCypressId () {
      return `button-${this.getSuspendStatus.toLowerCase()}`
    },

    systemFields () {
      return getSystemFields(this.role)
    },
  },
}
</script>
