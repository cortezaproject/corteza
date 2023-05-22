<template>
  <b-modal
    v-model="showModal"
    :title="alterationModalTitle"
    size="lg"
    body-class="p-0 border-top-0"
    header-class="p-3 pb-0 border-bottom-0"
    @change="$emit('change', $event)"
  >
    <b-table-simple hover small caption-top responsive class="px-3">
      <b-thead>
        <b-tr>
          <b-th>Operation</b-th>
          <b-th>Change</b-th>
          <b-th></b-th>
        </b-tr>
      </b-thead>
      <b-tbody>
        <template
          v-for="batch in alterationsBatched"
        >
          <b-tr
            v-for="a in batch"
            :key="a.alterationID"
          >
            <b-td>
              {{ a.alterationID }}
              {{ a.kind }}

              <template v-if="a.error">
                <br />
                <b-badge variant="danger">
                  {{ a.error || '' }}
                </b-badge>
              </template>
            </b-td>
            <b-td>
              <pre class="m-0">{{ a.params }}</pre>
            </b-td>

            <b-td>
              <c-input-confirm
                v-if="canComplete(a)"
                class="ml-1"
                @click.stop
                @confirmed="onComplete(a)"
              >
                Complete
              </c-input-confirm>
              <c-input-confirm
                v-if="canDismiss(a)"
                class="ml-1"
                @click.stop
                @confirmed="onDismiss(a)"
              >
                Dismiss
              </c-input-confirm>
              <b-badge
                v-else-if="!canDismiss(a) && !a.completedAt && a.dependsOn"
                variant="secondary"
              >
                Waiting for {{ a.dependsOn }}
              </b-badge>
              <b-badge
                v-else-if="a.completedAt"
                variant="success"
              >
                Completed
              </b-badge>
            </b-td>
          </b-tr>

        </template>
      </b-tbody>
    </b-table-simple>

    <b-container
      slot="modal-footer"
    >
      <b-row
        no-gutters
        class="align-items-center"
      >
        <b-col
          class="d-flex justify-content-start"
        >
          <b-button
            variant="light"
            data-test-id="button-save"
            size="sm"
            @click="$emit('save')"
          >
            Close
          </b-button>
        </b-col>
        <b-col
          class="d-flex justify-content-end"
        >
          <c-input-confirm
            class="mr-1"
            variant="primary"
            :borderless="false"
            @confirmed="onComplete()"
          >
            Complete automatically
          </c-input-confirm>
        </b-col>
      </b-row>
    </b-container>
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

      alterations: [{
        alterationID: '338079440671000001',
        batchID: '338079440672000001',
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
      }, {
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
        error: 'Something went wrong :)',
        createdAt: '2023-05-22T08:39:35.341Z',
        createdBy: '338079375675130271',
        updatedAt: undefined,
        updatedBy: undefined,
        deletedAt: undefined,
        deletedBy: undefined,
        completedAt: undefined,
        completedBy: undefined,
      }],
    }
  },

  computed: {
    alterationModalTitle () {
      const { handle } = this.module
      return handle ? this.$t('edit.alterationSettings.specificTitle', { handle }) : this.$t('edit.alterationSettings.title')
    },

    alterationsBatched () {
      let out = []
      let batchIndex = {}

      for (const a of this.alterations || []) {
        if (batchIndex[a.batchID]) {
          out[batchIndex[a.batchID]].push(a)
          continue
        }

        batchIndex[a.batchID] = out.length
        out.push([a])
      }

      return out
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
      if (!alteration) {
        alteration = this.alterations
          .filter(a => !a.completedAt)
          .map(a => a.alterationID)
      }

      // @todo here, we'll be pinging an API which will take a little bit to process
      await this.delay(1000)
      console.log('done on dismiss')
    },

    async onComplete (alteration = undefined) {
      if (!alteration) {
        alteration = this.alterations
          .filter(a => !a.completedAt)
          .map(a => a.alterationID)
      }

      // @todo here, we'll be pinging an API which will take a little bit to process
      await this.delay(1000)
      console.log('done on complete')
    },

    canDismiss (alteration) {
      if (alteration.completedAt) {
        return false
      }

      if (alteration.dependsOn) {
        const dep = this.alterations.find(a => a.alterationID === alteration.dependsOn)
        if (dep && !dep.completedAt) {
          return false
        }
      }

      return true
    },

    canComplete (alteration) {
      return this.canDismiss(alteration)
    },

    delay (t, val) {
      return new Promise(resolve => setTimeout(resolve, t, val));
    }
  },
}
</script>
