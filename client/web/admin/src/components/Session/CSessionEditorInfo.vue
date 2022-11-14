<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      v-if="session && session.sessionID"
    >
      <b-form-group
        :label="$t('id')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          :value="session.sessionID"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="session.workflowID"
        :label="$t('workflowID')"
        label-cols="2"
      >
        <b-form-input
          :value="session.workflowID"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        :label="$t('status')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          :value="session.status"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="session.error"
        :label="$t('error')"
        label-cols="2"
      >
        <b-form-input
          :value="session.error"
          plaintext
          disabled
          class="text-danger font-italic"
        />
      </b-form-group>

      <b-form-group
        v-if="session.resourceType"
        :label="$t('resourceType')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          :value="session.resourceType"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="session.eventType"
        :label="$t('eventType')"
        label-cols="2"
      >
        <b-form-input
          :value="session.eventType"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="session.workflowID"
        label-cols="2"
      >
        <b-button
          variant="light"
          :to="{ name: 'automation.workflow.edit', params: { workflowID: session.workflowID } }"
        >
          {{ $t('openWorkflow') }}
        </b-button>
      </b-form-group>

      <b-form-group
        v-if="!session.completedAt"
        label-cols="2"
      >
        <b-button
          variant="danger"
          :disabled="processing"
          @click="$emit('cancel')"
        >
          {{ $t('cancel') }}
        </b-button>
      </b-form-group>

      <b-form-group
        v-if="session.createdBy"
        :label="$t('createdByUserID')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          :value="session.createdBy"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="createdByUserName"
        :label="$t('createdByUserName')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          :value="createdByUserName"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="session.completedAt"
        :label="$t('completedAt')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          :value="session.completedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="session.createdAt"
        :label="$t('createdAt')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          :value="session.createdAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="session.deletedAt"
        :label="$t('deletedAt')"
        label-cols="2"
      >
        <b-form-input
          :value="session.deletedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>
    </b-form>

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>
  </b-card>
</template>

<script>
export default {
  name: 'CSessionEditorInfo',

  i18nOptions: {
    namespaces: 'automation.sessions',
    keyPrefix: 'editor.info',
  },

  props: {
    session: {
      type: Object,
      required: true,
    },

    user: {
      type: Object,
      required: true,
    },

    processing: {
      type: Boolean,
      value: false,
    },
  },

  computed: {
    createdByUserName () {
      const { userID, name, username, email } = this.user
      return name || username || email || `<@${userID}>`
    },
  },
}
</script>
