<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div class="my-2">
      <b-nav
        :tabs="options.display.appearance === 'tabs'"
        :pills="options.display.appearance === 'pills'"
        :small="options.display.appearance === 'small'"
        :fill="options.display.fillJustify === 'fill'"
        :justify="options.display.fillJustify === 'justify'"
        :align="options.display.alignment"
      >
        <b-nav-item
          v-for="(list, index) in options.lists"
          :key="`lists-${index}`"
          :to="processLink(list)"
          :href="processHrefLink(list)"
          :disabled="list.options.disabled"
          :variant="list.options.bgVariant"
          :style="{ order: index, color: list.options.textColor }"
          :link-attrs="{ style: `color: ${list.options.textColor}` }"
          :target="list.options.newWindow ? '_blank' : '_self'"
        >
          <template v-if="list.navigationType === 'text'">
            <span>
              {{ list.options.itemOption.text }}
            </span>
          </template>
          <template v-if="list.navigationType === 'url'">
            <span>
              {{ list.options.itemOption.text }}
            </span>
          </template>
          <template v-if="list.navigationType === 'compose'">
            <span>
              {{ list.options.itemOption.text }}
            </span>
          </template>
        </b-nav-item>

        <b-nav-item-dropdown
          v-for="(list, index) in options.lists"
          v-show="list.navigationType === 'dropdown'"
          :key="`lists-dropdown-${index}`"
          :text="list.options.dropdownText"
          :variant="list.options.bgVariant"
          right
          :style="{ order: index }"
        >
          <b-dropdown-item
            v-for="(dropdown, dIndex) in list.options.dropdown"
            :key="`dropdown-${dIndex}`"
            :to="dropdown.url"
            :disabled="list.options.disabled"
          >
            {{ dropdown.text }}
          </b-dropdown-item>
        </b-nav-item-dropdown>
      </b-nav>
    </div>
  </wrap>
</template>
<script>
import base from './base'

export default {
  extends: base,

  methods: {
    processLink (list) {
      if (list.navigationType === 'compose') {
        return list.options.itemOption.referenceId
      }

      return ''
    },

    processHrefLink (list) {
      if (list.navigationType === 'url') {
        return list.options.itemOption.url
      }

      return ''
    },
  },
}
</script>
