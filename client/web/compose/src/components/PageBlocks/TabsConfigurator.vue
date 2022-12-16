<template>
  <b-tab title="Tabs">
    <div>
      <h5 class="text-primary">
        {{ $t('tabs.displayTitle') }}
      </h5>

      <b-row
        class="mb-3 mt-3 ml-0 mr-0 justify-content-between"
        no-gutters
      >
        <b-form-group label="Appearance">
          <b-form-radio-group
            v-model="block.options.style.appearance"
            buttons
            button-variant="outline-primary"
            size="sm"
            :options="style.appearance"
          />
        </b-form-group>

        <b-form-group label="Alignment">
          <b-form-radio-group
            v-model="block.options.style.alignment"
            buttons
            button-variant="outline-primary"
            size="sm"
            :options="style.alignment"
          />
        </b-form-group>
        <b-form-group label="Fill or Justify">
          <b-form-radio-group
            v-model="block.options.style.fillJustify"
            buttons
            button-variant="outline-primary"
            size="sm"
            :options="style.fillJustify"
          />
        </b-form-group>
      </b-row>
    </div>

    <div
      class="d-flex"
    >
      <h5
        class="font-weight-light m-0 p-0 text-primary"
      >
        {{ $t('tabs.title') }}
      </h5>

      <b-button
        variant="link"
        size="md"
        :title="shouldDisableAdd ? $t('tabs.tooltip.tabCondition') : $t('tabs.tooltip.addTab')"
        :disabled="shouldDisableAdd"
        class="p-0 ml-3 text-decoration-none"
        @click="addTab"
      >
        {{ $t('tabs.addTab') }}
      </b-button>
    </div>

    <b-table-simple
      v-if="block.options.tabs.length"
      borderless
      small
    >
      <b-thead>
        <tr>
          <th />

          <th
            class="text-primary"
          >
            {{ $t('tabs.table.columns.title.label') }}
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
          <b-td class="handle align-middle">
            <font-awesome-icon
              :icon="['fas', 'bars']"
              class="grab m-0 text-light p-0"
            />
          </b-td>

          <b-td
            class="align-middle"
            style="width: 50%"
          >
            <b-form-input
              v-model="tab.title"
              :title="$t('tabs.tooltip.title')"
              :disabled="!tab.blockID"
              :placeholder="$t('tabs.form.title')"
            />
          </b-td>

          <b-td
            class="align-middle"
            style="width: 50%"
          >
            <div
              class="d-flex"
            >
              <vue-select
                v-model="tab.blockID"
                :title="$t('tabs.tooltip.selectBlock')"
                :options="options"
                :placeholder="$t('tabs.form.placeholder')"
                :selectable="option => isSelectable(option)"
                class="block-selector bg-white m-0"
                append-to-body
                style="min-width: 95%;"
                :reduce="option => option.value"
              >
                <template #list-footer>
                  <b-button
                    id="CreateBlockSelectorTab"
                    variant="link"
                    size="sm"
                    :title="$t('tabs.tooltip.newBlock')"
                    class="text-decoration-none"
                    block
                    @click="showBlockSelector(index)"
                  >
                    {{ $t('tabs.addTab') }}
                  </b-button>
                </template>
              </vue-select>

              <b-button
                id="popover-edit"
                size="sm"
                :disabled="!tab.blockID"
                :title="$t('tabs.tooltip.edit')"
                variant="light"
                @click="editBlock(tab.blockID)"
              >
                <font-awesome-icon
                  :icon="['far', 'edit']"
                />
              </b-button>
            </div>
          </b-td>

          <td
            class="text-right align-middle pr-2"
            style="min-width: 100px;"
          >
            <c-input-confirm
              :title="$t('tabs.tooltip.delete')"
              @confirmed="deleteTab(index)"
            />
          </td>
        </tr>
      </draggable>
    </b-table-simple>

    <div
      v-else
      class="text-center pt-5 pb-5"
    >
      <p>
        {{ $t('tabs.noTabs') }}
      </p>
    </div>

    <b-modal
      id="createBlockSelectorTab"
      size="lg"
      scrollable
      hide-footer
      :title="$t('tabs.newBlockModal')"
    >
      <new-block-selector
        :record-page="!!module"
        :disable-kind="['Tabs']"
        @select="addBlock"
      />
    </b-modal>
  </b-tab>
</template>

<script>
import base from './base'
import draggable from 'vuedraggable'
import { VueSelect } from 'vue-select'
import { fetchID } from 'corteza-webapp-compose/src/lib/tabs.js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'TabConfigurator',

  components: {
    draggable,
    VueSelect,
    //  Importing like this because configurator is recursive
    NewBlockSelector: () => import('corteza-webapp-compose/src/components/Admin/Page/Builder/Selector'),
  },

  extends: base,

  data () {
    return {
      activeIndex: null,
      style: {
        appearance: [
          { text: this.$t('tabs.style.appearance.tabs'), value: 'tabs', disabled: false },
          { text: this.$t('tabs.style.appearance.pills'), value: 'pills', disabled: false },
          { text: this.$t('tabs.style.appearance.small'), value: 'small', disabled: false },
        ],

        alignment: [
          { text: this.$t('tabs.style.alignment.left'), value: 'left', disabled: false },
          { text: this.$t('tabs.style.alignment.center'), value: 'center', disabled: false },
          { text: this.$t('tabs.style.alignment.right'), value: 'right', disabled: false },
        ],

        fillJustify: [
          { text: this.$t('tabs.style.fillJustify.fill'), value: 'fill', disabled: false },
          { text: this.$t('tabs.style.fillJustify.justified'), value: 'justified', disabled: false },
          { text: this.$t('tabs.style.fillJustify.none'), value: 'none', disabled: false },
        ],
      },
      untabbedBlock: [],
    }
  },

  computed: {

    options () {
      return this.page.blocks.filter(b => b.kind !== 'Tabs').map((b, i) => {
        // block title is going to look ugly till you save the page. Inevitable.
        return { value: fetchID(b), label: b.title || `Block-${b.kind}` }
      })
    },

    shouldDisableAdd () {
      return this.page.blocks.find(b => fetchID(b) === fetchID(this.block)) === undefined
    },
  },

  created () {
    this.$root.$on('builder-createRequestFulfilled', this.createRequestFulfilled)
  },

  destroyed () {
    this.$root.$off('builder-createRequestFulfilled', this.createRequestFulfilled)
  },

  methods: {
    createRequestFulfilled ({ tab }) {
      if (tab) {
        this.updateTab(tab, this.activeIndex)
      }
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
      block.meta.tabbed = true
      this.$root.$emit('tab-createRequest', block)
    },

    updateTab (tab, index) {
      this.block.options.tabs.splice(index, 1, tab)
    },

    deleteTab (tabIndex) {
      this.block.options.tabs.splice(tabIndex, 1)
    },
  },
}
</script>

<style lang="scss">
.block-selector {
  .vs__selected-options {
    flex-wrap: nowrap;
  }

  .vs__selected {
    max-width: 200px;
    display: inline-block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.vs__dropdown-menu .vs__dropdown-option {
  text-overflow: ellipsis;
  overflow: hidden !important;
}

</style>
