<template>
  <b-tab :title="$t('group.label')">
    <div
      class="d-flex align-items-center mb-2"
    >
      <h5 class="m-0">
        {{ $t('group.blocks.title') }}
      </h5>
    </div>

    <c-form-table-wrapper
      :labels="{
        addButton: $t('general:label.add')
      }"
      @add-item="showBlockSelector"
    >
      <b-table-simple
        responsive
        borderless
        small
      >
        <b-thead>
          <tr>
            <th scope="col" />
            <th
              class="text-primary"
              scope="col"
            >
              {{ $t('group.table.columns.block.label') }}
            </th>
            <th scope="col" />
          </tr>
        </b-thead>

        <draggable
          v-model="block.options.blocks"
          handle=".handle"
          tag="b-tbody"
        >
          <tr
            v-for="(b, index) in block.options.blocks"
            :key="index"
          >
            <b-td class="handle align-middle pr-2">
              <font-awesome-icon
                :icon="['fas', 'bars']"
                class="grab m-0 text-secondary p-0"
              />
            </b-td>

            <b-td class="align-middle">
              <b-input-group class="d-flex flex-nowrap w-100">
                <c-input-select
                  v-model="b.blockID"
                  :options="blockOptions"
                  :placeholder="$t('group.placeholder.block')"
                  :get-option-label="getBlockLabel"
                  :get-option-key="getOptionKey"
                  :selectable="option => isSelectable(option)"
                  :reduce="option => option.value"
                />

                <b-input-group-append>
                  <b-button
                    v-if="b.blockID"
                    v-b-tooltip.noninteractive.hover="{ title: $t('group.tooltip.edit'), boundary: 'body' }"
                    size="sm"
                    variant="extra-light"
                    class="d-flex align-items-center justify-content-center"
                    style="width: 40px;"
                    @click="editBlock(b.blockID)"
                  >
                    <font-awesome-icon :icon="['far', 'edit']" />
                  </b-button>
                </b-input-group-append>
              </b-input-group>
            </b-td>

            <td
              class="text-center align-middle"
              style="min-width: 80px;"
            >
              <c-input-confirm
                :tooltip="$t('group.tooltip.delete')"
                show-icon
                @confirmed="deleteBlock(index)"
              />
            </td>
          </tr>
        </draggable>
      </b-table-simple>
    </c-form-table-wrapper>

    <b-modal
      id="createBlockSelectorGroup"
      size="lg"
      scrollable
      hide-footer
      no-fade
      :title="$t('group.newBlockModal')"
    >
      <new-block-selector
        :record-page="!!module"
        :disabled-kinds="['Tabs', 'Group']"
        @select="addBlock"
      />
    </b-modal>
  </b-tab>
</template>

<script>
import base from './base'
import draggable from 'vuedraggable'
import { fetchID } from 'corteza-webapp-compose/src/lib/block'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'GroupConfigurator',

  components: {
    draggable,
    NewBlockSelector: () => import('corteza-webapp-compose/src/components/Admin/Page/Builder/Selector'),
  },

  extends: base,

  computed: {
    blockOptions () {
      return [
        ...this.page.blocks.filter(({ blockID, kind }) => !['Tabs', 'Group'].includes(kind) && !this.blocks.some(b => b.blockID === blockID) && this.block.options.blocks.some(b => b.blockID === blockID)),
        ...this.blocks.filter(({ kind }) => !['Tabs', 'Group'].includes(kind)),
      ].map(b => ({ ...b, value: fetchID(b) }))
    },
  },

  mounted () {
    this.$root.$on('builder-createRequestFulfilled', this.createRequestFulfilled)
  },

  beforeDestroy () {
    this.$root.$off('builder-createRequestFulfilled', this.createRequestFulfilled)
  },

  methods: {
    createRequestFulfilled (block) {
      this.block.options.blocks.push({
        blockID: fetchID(block),
        xywh: [0, 0, 12, 10], // Default size inside group
      })
    },

    showBlockSelector () {
      this.$bvModal.show('createBlockSelectorGroup')
    },

    addBlock (block) {
      this.$bvModal.hide('createBlockSelectorGroup')
      block.meta.hidden = true
      this.$root.$emit('tab-createRequest', block) // Reusing tab-createRequest as it's handled by Builder
    },

    editBlock (blockID) {
      this.$root.$emit('tab-editRequest', blockID)
    },

    deleteBlock (index) {
      this.block.options.blocks.splice(index, 1)
    },

    getBlockLabel ({ title, kind }) {
      return title || kind
    },

    getOptionKey (block) {
      return fetchID(block)
    },

    isSelectable (option) {
      return !this.block.options.blocks.some(b => b.blockID === option.value)
    },
  },
}
</script>
