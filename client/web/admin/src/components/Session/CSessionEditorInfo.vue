<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap gap-1"
  >
    <b-form
      v-if="session && session.sessionID"
    >
      <b-row>
        <b-col
          v-if="session.workflowID"
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('workflowID')"
            label-class="text-primary"
          >
            {{ session.workflowID }}
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('status')"
            label-class="text-primary"
            class="mb-0"
          >
            {{ session.status }}
          </b-form-group>
        </b-col>

        <b-col
          v-if="session.error"
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('error')"
            label-class="text-primary"
          >
            {{ session.error }}
          </b-form-group>
        </b-col>

        <b-col
          v-if="session.resourceType"
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('resourceType')"
            label-class="text-primary"
          >
            {{ session.resourceType }}
          </b-form-group>
        </b-col>

        <b-col
          v-if="session.eventType"
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('eventType')"
            label-class="text-primary"
          >
            {{ session.eventType }}
          </b-form-group>
        </b-col>

        <b-col
          v-if="!session.completedAt"
          cols="12"
        >
          <b-form-group
            label-class="text-primary"
          >
            <b-button
              variant="danger"
              :disabled="processing"
              @click="$emit('cancel')"
            >
              {{ $t('cancel') }}
            </b-button>
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :id="session.sessionID"
        :resource="session"
      >
        <template #custom-field>
          <b-col
            v-if="createdByUserName"
            cols="12"
            lg="4"
          >
            <b-form-group
              :label="$t('createdByUserName')"
              label-class="text-primary"
            >
              {{ createdByUserName }}
            </b-form-group>
          </b-col>
        </template>
      </c-system-fields>
    </b-form>

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <template
      v-if="session.workflowID"
      #footer
    >
      <b-button
        variant="light"
        :to="{ name: 'automation.workflow.edit', params: { workflowID: session.workflowID } }"
      >
        {{ $t('openWorkflow') }}
      </b-button>
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
