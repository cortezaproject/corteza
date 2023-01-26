<template>
  <div>
    <b-tab :title="$t('navigation.label')">
      <div class="mb-3">
        <h5 class="text-primary">
          {{ $t("navigation.displayOptions") }}
        </h5>

        <b-row class="justify-content-between">
          <b-col
            cols="12"
            sm="4"
            class="mb-2 mb-sm-0"
          >
            <b-form-group
              horizontal
              :label="$t('navigation.appearance')"
            >
              <b-form-radio-group
                v-model="options.display.appearance"
                buttons
                button-variant="outline-secondary"
                :options="appearanceOptions"
              />
            </b-form-group>
          </b-col>

          <b-col
            sm="4"
            class="mb-2 mb-sm-0"
          >
            <b-form-group
              horizontal
              :label="$t('navigation.fillJustify')"
            >
              <b-form-radio-group
                v-model="options.display.fillJustify"
                buttons
                button-variant="outline-secondary"

                :options="fillJustifyOptions"
              />
            </b-form-group>
          </b-col>

          <b-col
            sm="4"
            class="mb-2 mb-sm-0"
          >
            <b-form-group
              horizontal
              :label="$t('navigation.alignment')"
            >
              <b-form-radio-group
                v-model="options.display.alignment"
                buttons
                button-variant="outline-secondary"
                :options="alignmentOptions"
              />
            </b-form-group>
          </b-col>
        </b-row>
      </div>

      <div class="mb-3 mt-2">
        <div class="d-flex align-items-center mb-4">
          <h5 class="text-primary mb-0">
            {{ $t("navigation.navigationItems") }}
          </h5>

          <b-button
            variant="link"
            class="d-flex align-items-center text-decoration-none"
            @click="addNavigationItem"
          >
            <font-awesome-icon
              :icon="['fas', 'plus']"
              size="sm"
              class="mr-1"
            />
            {{ $t("navigation.add") }}
          </b-button>
        </div>

        <div class="mt-3">
          <draggable
            v-model="block.options.navigationItems"
            group="sort"
            handle=".grab"
          >
            <div
              v-for="(item, index) in block.options.navigationItems"
              :key="index"
            >
              <hr v-if="index">

              <b-table-simple
                borderless
              >
                <thead>
                  <tr>
                    <th style="width: auto;" />
                    <th>
                      {{ $t("navigation.type") }}
                    </th>
                    <th>
                      {{ $t("navigation.color") }}
                    </th>
                    <th>
                      {{ $t("navigation.background") }}
                    </th>
                    <th class="text-center">
                      {{ $t("navigation.enabled") }}
                    </th>
                    <th style="width: auto; min-width: 100px;" />
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td class="align-middle">
                      <font-awesome-icon
                        :icon="['fas', 'bars']"
                        class="grab text-grey"
                      />
                    </td>
                    <td class="align-middle">
                      <b-form-select
                        v-model="item.type"
                        :options="navigationItemTypes"
                      />
                    </td>
                    <td class="align-middle">
                      <b-form-input
                        v-model="item.options.textColor"
                        type="color"
                        debounce="300"
                        class="w-100"
                      />
                    </td>
                    <td class="align-middle">
                      <b-form-input
                        v-model="item.options.backgroundColor"
                        type="color"
                        debounce="300"
                        class="w-100"
                      />
                    </td>
                    <td class="align-middle text-center">
                      <b-form-checkbox
                        v-model="item.options.enabled"
                        switch
                        class="mb-0"
                      />
                    </td>
                    <td class="align-middle">
                      <c-input-confirm
                        button-class="px-2"
                        size="md"
                        @confirmed="options.navigationItems.splice(index, 1)"
                      />
                    </td>
                  </tr>

                  <component
                    :is="item.type"
                    :item="item"
                    :namespace="namespace"
                    @update="(value) => item = value"
                  />
                </tbody>
              </b-table-simple>
            </div>
          </draggable>
        </div>
      </div>
    </b-tab>
  </div>
</template>

<script>
import base from '../base'
import Draggable from 'vuedraggable'
import { compose } from '@cortezaproject/corteza-js'
import Text from './NavTypes/Text.vue'
import Url from './NavTypes/Url.vue'
import Compose from './NavTypes/ComposePage.vue'
import Dropdown from './NavTypes/Dropdown.vue'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    Draggable,
    TextSection: Text,
    Url,
    Compose,
    Dropdown,
  },

  extends: base,

  data () {
    return {
      appearanceOptions: [
        { value: 'tabs', text: this.$t('navigation.tabs') },
        { value: 'pills', text: this.$t('navigation.pills') },
        { value: 'small', text: this.$t('navigation.small') },
      ],

      alignmentOptions: [
        { value: 'left', text: this.$t('navigation.left') },
        { value: 'center', text: this.$t('navigation.center') },
        { value: 'right', text: this.$t('navigation.right') },
      ],

      fillJustifyOptions: [
        { value: 'fill', text: this.$t('navigation.fill') },
        { value: 'justify', text: this.$t('navigation.justify') },
      ],

      backgroundColors: [
        { value: 'primary', text: this.$t('navigation.primary') },
        { value: 'secondary', text: this.$t('navigation.secondary') },
        { value: 'success', text: this.$t('navigation.success') },
        { value: 'warning', text: this.$t('navigation.warning') },
        { value: 'danger', text: this.$t('navigation.danger') },
        { value: 'info', text: this.$t('navigation.info') },
      ],

      navigationItemTypes: [
        { value: 'url', text: this.$t('navigation.url') },
        { value: 'compose', text: this.$t('navigation.composePage') },
        { value: 'dropdown', text: this.$t('navigation.dropdown') },
        { value: 'text-section', text: this.$t('navigation.text') },
      ],
    }
  },

  methods: {
    addNavigationItem () {
      this.block.options.navigationItems.push(
        compose.PageBlockNavigation.makeNavigationItem({
          type: 'compose',
          options: {
            backgroundColor: '#ffffff',
            item: {
              align: 'bottom',
              target: 'sameTab',
              dropdown: {
                label: '',
                items: [],
              },
            },
          },
        }),
      )
    },
  },
}
</script>

<style lang="scss" scoped>
th {
  width: 25%;
}

th,
td {
  padding-left: 15px;
  padding-right: 15px;
}
</style>
