<template>
  <div>
    <div class="mb-5">
      <h5 class="text-primary">
        {{ $t('navigation.displayOptions') }}
      </h5>

      <b-row class="justify-content-between">
        <b-col
          cols="6"
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
              size="sm"
              name="appearance"
              :options="appearanceOptions"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="6"
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
              size="sm"
              name="alignment"
              :options="alignmentOptions"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="6"
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
              size="sm"
              name="fill-justify"
              :options="fillJustifyOptions"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </div>
    <div class="mb-5">
      <h5 class="text-primary">
        {{ $t('navigation.listItems') }}
      </h5>

      <div class="mt-3">
        <draggable
          v-model="block.options.lists"
          group="sort"
          handle=".grab"
        >
          <b-form-row
            v-for="(column, index) in block.options.lists"
            :key="index"
            class="mb-1 border-bottom mt-3"
          >
            <b-col
              cols="1"
              class="d-flex mt-4 justify-content-center"
            >
              <font-awesome-icon
                :icon="['fas', 'bars']"
                class="grab text-grey"
              />
            </b-col>

            <b-col cols="2">
              <b-form-group :label="$t('navigation.type')">
                <b-form-select
                  v-model="column.navigationType"
                  :options="navigationItemTypes"
                  size="sm"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="7"
              class="mb-5"
            >
              <template>
                <h6>{{ $t('navigation.displayOptions') }}</h6>

                <b-row align-v="center">
                  <b-col>
                    <b-form-group>
                      <b-form-checkbox
                        v-model="column.options.disabled"
                        name="disabled-button"
                        switch
                      >
                        {{ $t('navigation.disabled') }}
                      </b-form-checkbox>
                    </b-form-group>
                  </b-col>
                  <b-col>
                    <b-form-group :label="$t('navigation.color')">
                      <b-form-input
                        v-model="column.options.textColor"
                        size="sm"
                        type="color"
                        style="width: 80px"
                      />
                    </b-form-group>
                  </b-col>
                  <b-col>
                    <b-form-group :label="$t('navigation.background')">
                      <b-form-select
                        v-model="column.options.bgVariant"
                        :options="bgVariants"
                        size="sm"
                        type="color"
                      />
                    </b-form-group>
                  </b-col>
                </b-row>
              </template>

              <template v-if="column.navigationType === 'text'">
                <text-type :column.sync="column" />
              </template>

              <template v-if="column.navigationType === 'url'">
                <url-type :column.sync="column" />
              </template>

              <template v-if="column.navigationType === 'dropdown'">
                <dropdown-type :column.sync="column" />
              </template>

              <template v-if="column.navigationType === 'compose'">
                <compose-type
                  :column.sync="column"
                  :namespace="namespace"
                />
              </template>
            </b-col>

            <b-col
              cols="2"
              class="d-flex mt-4 justify-content-around"
            >
              <c-input-confirm
                variant="link"
                size="lg"
                button-class="text-dark px-0"
                @confirmed="options.lists.splice(index, 1)"
              />
            </b-col>
          </b-form-row>
        </draggable>

        <div class="d-flex align-items-center mt-4">
          <b-button
            class="d-flex align-items-center text-decoration-none"
            @click="handleAddButton"
          >
            <font-awesome-icon
              :icon="['fas', 'plus']"
              size="sm"
              class="mr-1"
            />
            {{ $t('navigation.add') }}
          </b-button>
        </div>
      </div>
    </div>

    <div class="my-5 py-5 py-sm-3" />

    <div class="preview bg-white position-absolute p-3">
      <h6 class="text-primary">
        {{ $t('navigation.preview.label') }}
      </h6>

      <b-nav pills>
        <b-nav-item active>
          {{ $t('navigation.preview.home') }}
        </b-nav-item>
        <b-nav-item>{{ $t('navigation.preview.about') }}</b-nav-item>
        <b-nav-item-dropdown
          id="preview-dropdown"
          text="Preview dropdown"
          toggle-class="nav-link-custom"
          right
        >
          <b-dropdown-item>{{ $t('navigation.preview.dropdown') }}</b-dropdown-item>
          <b-dropdown-item>{{ $t('navigation.preview.dropdown') }}</b-dropdown-item>
        </b-nav-item-dropdown>
      </b-nav>
    </div>
  </div>
</template>
<script>
import base from '../base'
import Draggable from 'vuedraggable'
import { compose } from '@cortezaproject/corteza-js'
import TextType from './NavTypes/TextType.vue'
import UrlType from './NavTypes/UrlType.vue'
import ComposeType from './NavTypes/ComposeType.vue'
import DropdownType from './NavTypes/DropdownType.vue'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: { Draggable, TextType, UrlType, ComposeType, DropdownType },

  extends: base,

  data () {
    return {
      appearanceOptions: [
        { value: 'tabs', text: this.$t('navigation.tabs') },
        { value: 'pills', text: this.$t('navigation.pills') },
        { value: 'small', text: this.$t('navigation.small') },
      ],

      alignmentOptions: [
        { value: 'center', text: this.$t('navigation.center') },
        { value: 'left', text: this.$t('navigation.left') },
        { value: 'right', text: this.$t('navigation.right') },
      ],

      fillJustifyOptions: [
        { value: 'fill', text: 'Fill' },
        { value: 'justify', text: 'Justify' },
      ],

      bgVariants: [
        { value: 'primary', text: this.$t('navigation.primary') },
        { value: 'secondary', text: this.$t('navigation.secondary') },
        { value: 'success', text: this.$t('navigation.success') },
        { value: 'warning', text: 'Warning' },
        { value: 'danger', text: this.$t('navigation.danger') },
        { value: 'info', text: this.$t('navigation.info') },
      ],

      navigationItemTypes: [
        { value: 'text', text: this.$t('navigation.text') },
        { value: 'url', text: this.$t('navigation.url') },
        { value: 'compose', text: this.$t('navigation.composePage') },
        { value: 'dropdown', text: this.$t('navigation.dropdown') },
      ],

      dropdownTypes: [
        { value: 'text', text: this.$t('navigation.text') },
        { value: 'multi', text: this.$t('navigation.multi') },
      ],

    }
  },

  methods: {
    handleAddButton () {
      this.block.options.lists.push(compose.PageBlockNavigation.makeListItem())
    },
  },
}
</script>

<style lang="scss" scoped>
.preview {
  bottom: 0;
  left: 0;
  z-index: 2;
  width: 100%;
  box-shadow: 0 -0.25rem 1rem rgb(0 0 0 / 15%);
}
</style>
