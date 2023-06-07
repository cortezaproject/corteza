<template>
  <b-modal
    v-model="showModal"
    title="Schema alterations"
    size="xl"
    body-class="p-0 border-top-0"
    header-class="p-3 pb-0 border-bottom-0"
    @change="$emit('change', $event)"
  >
    <b-table-simple
      borderless
      sticky-header
      responsive
      head-variant="light"
      style="max-height: 75vh;"
    >
      <b-thead>
        <b-tr>
          <b-th
            class="text-primary"
          >
            Operation
          </b-th>

          <b-th
            class="text-primary"
          >
            Change
          </b-th>

          <b-th
            class="text-primary text-center"
          >
            Status
          </b-th>
          <b-th style="min-width: 200px;" />
        </b-tr>
      </b-thead>

      <b-tbody>
        <b-tr
          v-for="a in sortedAlterations"
          :key="a.alterationID"
          class="border-top"
          :class="{ 'bg-gray': a.alterationID === dependOnHover }"
          @mouseover="dependOnHover = a.dependsOn"
          @mouseleave="dependOnHover = undefined"
        >
          <b-td>
            {{ a.alterationID }}
          </b-td>

          <b-td>
            <pre class="m-0">{{ a.params }}</pre>
          </b-td>

          <b-td class="text-center align-top">
            <b-badge
              v-if="a.error"
              variant="danger"
            >
              {{ a.error || '' }}
            </b-badge>

            <b-badge
              v-else-if="a.completedAt"
              variant="success"
            >
              Resolved
            </b-badge>

            <b-badge
              v-else-if="a.dependsOn"
              variant="secondary"
            >
              Waiting for {{ a.dependsOn }}
            </b-badge>
          </b-td>

          <b-td class="text-right">
            <b-spinner
              v-if="a.processing"
              variant="primary"
              small
            />

            <template v-else>
              <c-input-confirm
                v-if="!a.completedAt"
                :disabled="!canResolve(a) || a.processing || processing"
                variant="primary"
                size="sm"
                class="mx-1"
                @click.stop
                @confirmed="onResolve(a)"
              >
                Resolve
              </c-input-confirm>

              <c-input-confirm
                v-if="!a.completedAt"
                :disabled="!canDismiss(a) || a.processing || processing"
                variant="light"
                size="sm"
                class="mx-1"
                @click.stop
                @confirmed="onDismiss(a)"
              >
                Dismiss
              </c-input-confirm>
            </template>
          </b-td>
        </b-tr>
      </b-tbody>
    </b-table-simple>

    <template #modal-footer>
      <b-button
        variant="link"
        size="sm"
        :disabled="processing"
        @click="$emit('cancel')"
      >
        Cancel
      </b-button>

      <c-input-confirm
        variant="primary"
        :disabled="processing"
        :borderless="false"
        @confirmed="onResolve()"
      >
        Resolve automatically
        <b-spinner
          v-if="processing"
          variant="white"
          small
        />
      </c-input-confirm>
    </template>
  </b-modal>
</template>
<script>
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'module',
  },

  props: {
    modal: {
      type: Boolean,
      required: false,
    },

    module: {
      type: compose.Module,
      required: true,
    },
  },

  data () {
    return {
      showModal: false,

      processing: false,

      dependOnHover: undefined,

      alterations: [
        {
          alterationID: '338079440671000001',
          batchID: '338079440672000001',
          dependsOn: undefined,
          processing: false,
          kind: 'a',
          params: {},
          error: undefined,
          createdAt: '2023-05-22T08:39:35.341Z',
          createdBy: '338079375675130271',
          updatedAt: undefined,
          updatedBy: undefined,
          deletedAt: undefined,
          deletedBy: undefined,
          completedAt: undefined,
          completedBy: undefined,
        }, {
          alterationID: '338079440671000002',
          batchID: '338079440672000002',
          dependsOn: undefined,
          kind: 'b',
          params: {},
          error: undefined,
          createdAt: '2023-05-22T08:39:35.341Z',
          createdBy: '338079375675130271',
          updatedAt: undefined,
          updatedBy: undefined,
          deletedAt: undefined,
          deletedBy: undefined,
          completedAt: '2023-05-22T08:39:35.341Z',
          completedBy: '338079375675130271',
        },
        {
          alterationID: '338079440671000003',
          batchID: '338079440672000003',
          dependsOn: undefined,
          kind: 'a',
          params: {},
          error: undefined,
          createdAt: '2023-05-22T08:39:35.341Z',
          createdBy: '338079375675130271',
          updatedAt: undefined,
          updatedBy: undefined,
          deletedAt: undefined,
          deletedBy: undefined,
          completedAt: undefined,
          completedBy: undefined,
        },
        {
          alterationID: '338079440671000004',
          batchID: '338079440672000003',
          dependsOn: '338079440671000003',
          kind: 'b',
          params: {},
          error: undefined,
          createdAt: '2023-05-22T08:39:35.341Z',
          createdBy: '338079375675130271',
          updatedAt: undefined,
          updatedBy: undefined,
          deletedAt: undefined,
          deletedBy: undefined,
          completedAt: undefined,
          completedBy: undefined,
        },
        {
          alterationID: '338079440671000005',
          batchID: '338079440672000005',
          dependsOn: undefined,
          kind: 'b',
          params: {},
          error: 'Something went wrong',
          createdAt: '2023-05-22T08:39:35.341Z',
          createdBy: '338079375675130271',
          updatedAt: undefined,
          updatedBy: undefined,
          deletedAt: undefined,
          deletedBy: undefined,
          completedAt: undefined,
          completedBy: undefined,
        },
      ],
    }
  },

  computed: {
    sortedAlterations () {
      return this.alterations.toSorted((a, b) => (a.batchID || '').localeCompare(b.batchID || '') || (a.dependsOn || '').localeCompare(b.dependsOn || ''))
    },
  },

  watch: {
    modal: {
      immediate: true,
      handler (show = false) {
        this.showModal = show
      },
    },
  },

  methods: {
    async onDismiss (alteration = undefined) {
      this.processing = true
      this.$set(alteration, 'processing', true)

      this.$ComposeAPI.dismissAlteration(alteration).then(() => {
        this.toastSuccess(this.$t('notification:module.alteration.dismiss.success'))
      }).catch(this.toastErrorHandler(this.$t('notification:alteration.dismiss.error')))
        .finally(() => {
          this.processing = false
          this.$set(alteration, 'processing', false)
        })
    },

    async onResolve (alteration = undefined) {
      this.processing = true
      this.$set(alteration, 'processing', true)

      this.$ComposeAPI.resolveAlteration(alteration).then(() => {
        this.toastSuccess(this.$t('notification:module.alteration.resolve.success'))
      }).catch(this.toastErrorHandler(this.$t('notification:alteration.resolve.error')))
        .finally(() => {
          this.processing = false
          this.$set(alteration, 'processing', false)
        })
    },

    canDismiss (alteration) {
      if (alteration.completedAt) {
        return false
      }

      if (alteration.dependsOn) {
        return this.alterations.some(a => a.alterationID === alteration.dependsOn && !a.completedAt)
      }

      return true
    },

    canResolve (alteration) {
      return this.canDismiss(alteration)
    },
  },
}
</script>
