<template>
  <b-tab :title="$t('tabs.label')">
    <h5>
      {{ $t('tabs.displayTitle') }}
    </h5>

    <b-row
      class="text-primary"
      no-gutters
    >
      <b-col
        cols="12"
        lg="4"
      >
        <b-form-group
          :label="$t('tabs.style.appearance.label')"
          label-class="text-primary"
        >
          <b-form-radio-group
            v-model="block.options.style.appearance"
            :options="style.appearance"
            buttons
            button-variant="outline-secondary"
            size="sm"
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        lg="4"
      >
        <b-form-group
          :label="$t('tabs.style.orientation.label')"
          label-class="text-primary"
        >
          <b-form-radio-group
            v-model="block.options.style.orientation"
            :options="style.orientation"
            buttons
            button-variant="outline-secondary"
            size="sm"
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        lg="4"
      >
        <b-form-group
          :label="$t('tabs.style.position.label')"
          label-class="text-primary"
        >
          <b-form-radio-group
            v-model="block.options.style.position"
            :options="style.position"
            buttons
            button-variant="outline-secondary"
            size="sm"
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        lg="4"
      >
        <b-form-group
          :label="$t('tabs.style.alignment.label')"
          label-class="text-primary"
        >
          <b-form-radio-group
            v-model="block.options.style.alignment"
            :options="style.alignment"
            buttons
            button-variant="outline-secondary"
            size="sm"
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        lg="4"
      >
        <b-form-group
          :label="$t('tabs.style.justify.label')"
          label-class="text-primary"
        >
          <b-form-radio-group
            v-model="block.options.style.justify"
            :options="style.justifyOptions"
            buttons
            button-variant="outline-secondary"
            size="sm"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <hr class="my-2">

    <div
      class="d-flex align-items-center mb-2"
    >
      <h5 class="m-0">
        {{ $t('tabs.title') }}
      </h5>
    </div>

    <c-form-table-wrapper
      :labels="{
        addButton: $t('general:label.add')
      }"
      @add-item="addTab"
    >
      <b-table-simple
        responsive="lg"
        borderless
        small
      >
        <b-thead>
          <tr>
            <th />

            <th class="d-flex align-items-center text-primary">
              {{ $t('tabs.table.columns.title.label') }}
              <c-hint
                :tooltip="$t('interpolationFootnote', ['${record.values.fieldName}', '${recordID}', '${ownerID}', '${userID}'])"
                class="d-block"
              />
            </th>

            <th
              class="text-primary"
            >
              {{ $t('tabs.table.columns.block.label') }}
            </th>

            <th />
          </tr>
        </b-thead>

        <draggable
          v-model="block.options.tabs"
          handle=".handle"
          tag="b-tbody"
        >
          <tr
            v-for="(tab, index) in block.options.tabs"
            :key="index"
          >
            <b-td class="handle align-middle pr-2">
              <font-awesome-icon
                :icon="['fas', 'bars']"
                class="grab m-0 text-secondary p-0"
              />
            </b-td>

            <b-td
              class="align-middle"
              style="width: 50%; min-width: 200px;"
            >
              <b-form-input
                v-model="tab.title"
              />
            </b-td>

            <b-td
              class="align-middle"
              style="width: 50%; min-width: 200px;"
            >
              <b-input-group class="d-flex flex-nowrap w-100">
                <c-input-select
                  v-model="tab.blockID"
                  :options="blockOptions"
                  :placeholder="$t('tabs.placeholder.block')"
                  :get-option-label="getBlockLabel"
                  :get-option-key="getOptionKey"
                  :selectable="option => isSelectable(option)"
                  :reduce="option => option.value"
                />

                <b-input-group-append>
                  <b-button
                    v-if="tab.blockID"
                    id="popover-edit"
                    v-b-tooltip.noninteractive.hover="{ title: $t('tabs.tooltip.edit'), container: '#body' }"
                    size="sm"
                    variant="extra-light"
                    class="d-flex align-items-center justify-content-center"
                    style="width: 40px;"
                    @click="editBlock(tab.blockID)"
                  >
                    <font-awesome-icon
                      :icon="['far', 'edit']"
                    />
                  </b-button>

                  <b-button
                    v-else
                    v-b-tooltip.noninteractive.hover="{ title: $t('tabs.tooltip.addBlock'), container: '#body' }"
                    size="sm"
                    variant="extra-light"
                    class="d-flex align-items-center justify-content-center"
                    style="width: 40px;"
                    @click="showBlockSelector(index)"
                  >
                    <font-awesome-icon
                      :icon="['fas', 'plus']"
                    />
                  </b-button>
                </b-input-group-append>
              </b-input-group>
            </b-td>

            <td
              class="text-center align-middle"
              style="min-width: 80px;"
            >
              <c-input-confirm
                :tooltip="$t('tabs.tooltip.delete')"
                show-icon
                @confirmed="deleteTab(index)"
              />
            </td>
          </tr>
        </draggable>
      </b-table-simple>
    </c-form-table-wrapper>

    <b-modal
      id="createBlockSelectorTab"
      size="lg"
      scrollable
      hide-footer
      no-fade
      :title="$t('tabs.newBlockModal')"
    >
      <new-block-selector
        :record-page="!!module"
        :disabled-kinds="['Tabs']"
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

  name: 'TabConfigurator',

  components: {
    draggable,
    //  Importing like this because configurator is recursive
    NewBlockSelector: () => import('corteza-webapp-compose/src/components/Admin/Page/Builder/Selector'),
  },

  extends: base,

  data () {
    return {
      activeIndex: null,
      style: {
        appearance: [
          { text: this.$t('tabs.style.appearance.tabs'), value: 'tabs' },
          { text: this.$t('tabs.style.appearance.pills'), value: 'pills' },
          { text: this.$t('tabs.style.appearance.small'), value: 'small' },
        ],

        alignment: [
          { text: this.$t('tabs.style.alignment.left'), value: 'left' },
          { text: this.$t('tabs.style.alignment.center'), value: 'center' },
          { text: this.$t('tabs.style.alignment.right'), value: 'right' },
        ],

        justifyOptions: [
          { text: this.$t('tabs.style.justify.justify'), value: 'justify' },
          { text: this.$t('tabs.style.justify.none'), value: 'none' },
        ],

        orientation: [
          { text: this.$t('tabs.style.orientation.horizontal'), value: 'horizontal' },
          { text: this.$t('tabs.style.orientation.vertical'), value: 'vertical' },
        ],

        position: [
          { text: this.$t('tabs.style.position.start'), value: 'start' },
          { text: this.$t('tabs.style.position.end'), value: 'end' },
        ],
      },
    }
  },

  computed: {
    blockOptions () {
      return [
        ...this.page.blocks.filter(({ blockID, kind }) => kind !== 'Tabs' && !this.blocks.some(b => b.blockID === blockID) && this.options.tabs.some(b => b.blockID === blockID)),
        ...this.blocks.filter(b => b.kind !== 'Tabs'),
      ].map(b => ({ ...b, value: fetchID(b) }))
    },
  },

  mounted () {
    this.$root.$on('builder-createRequestFulfilled', this.createRequestFulfilled)
  },

  beforeDestroy () {
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    createRequestFulfilled (tab) {
      const { title = '' } = this.block.options.tabs[this.activeIndex] || {}
      tab.title = title
      this.updateTab(tab, this.activeIndex)
    },

    addTab () {
      this.block.options.tabs.push({
        blockID: null,
        title: undefined,
      })
    },

    isSelectable (option) {
      return !this.block.options.tabs.some(t => t.blockID === option.value)
    },

    showBlockSelector (index) {
      this.$bvModal.show('createBlockSelectorTab')
      this.activeIndex = index
    },

    editBlock (blockID = undefined) {
      this.$root.$emit('tab-editRequest', blockID)
    },

    addBlock (block) {
      this.$bvModal.hide('createBlockSelectorTab')
      block.meta.hidden = true
      this.$root.$emit('tab-createRequest', block)
    },

    updateTab (tab, index) {
      this.block.options.tabs.splice(index, 1, tab)
    },

    deleteTab (tabIndex) {
      this.block.options.tabs.splice(tabIndex, 1)
    },

    getBlockLabel ({ title, kind }) {
      return title || kind
    },

    getOptionKey (block) {
      return fetchID(block)
    },

    setDefaultValues () {
      this.activeIndex = null
      this.style = {}
    },

    destroyEvents () {
      this.$root.$off('builder-createRequestFulfilled', this.createRequestFulfilled)
    },
  },
}
</script>
