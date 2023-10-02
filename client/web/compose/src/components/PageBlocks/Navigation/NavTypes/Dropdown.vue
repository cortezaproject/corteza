<template>
  <tr>
    <td />

    <td
      colspan="5"
      class="p-0"
    >
      <div class="d-flex">
        <th style="min-width: 200px;">
          <b-form-group
            :label="$t('navigation.fieldLabel')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="options.item.dropdown.label"
              type="text"
            />
          </b-form-group>
        </th>

        <th style="min-width: 200px;">
          <b-form-group
            :label="$t('navigation.drop')"
            horizontal
            label-class="text-primary"
          >
            <b-form-radio-group
              v-model="options.item.align"
              buttons
              button-variant="outline-primary"
              size="sm"
              :options="aligns"
            />
          </b-form-group>
        </th>
      </div>

      <div class="mb-4 mb-3 px-3">
        <h6 class="text-primary mb-0">
          {{ $t("navigation.dropdownItems") }}
        </h6>
      </div>

      <div class="px-3">
        <table
          v-if="options.item.dropdown.items.length > 0"
          class="dd-items table table-sm table-borderless table-responsive-lg"
        >
          <tr>
            <th style="min-width: 200px;">
              {{ $t("navigation.text") }}
            </th>
            <th style="min-width: 200px;">
              {{ $t("navigation.url") }}
            </th>
            <th style="min-width: 200px;">
              {{ $t('navigation.openIn') }}
            </th>
            <th
              class="text-center"
              style="width: 50px; min-width: 50px;"
            >
              {{ $t("navigation.delimiter") }}
            </th>
          </tr>

          <tr
            v-for="(item, dropIndex) in options.item.dropdown.items"
            :key="`drop-${dropIndex}`"
          >
            <td>
              <b-form-group class="mb-0">
                <b-form-input
                  v-model="item.label"
                  type="text"
                />
              </b-form-group>
            </td>

            <td>
              <b-form-group class="mb-0">
                <b-form-input
                  v-model="item.url"
                  type="text"
                />
              </b-form-group>
            </td>

            <td class="align-middle text-center">
              <b-form-group class="mb-0">
                <b-form-select
                  v-model="item.target"
                  :options="targetOptions"
                />
              </b-form-group>
            </td>

            <td class="align-middle text-center">
              <b-form-group class="d-flex align-items-center justify-content-center mb-0">
                <c-input-checkbox
                  v-model="item.delimiter"
                  switch
                  size="sm"
                />
              </b-form-group>
            </td>

            <td class="align-middle text-center">
              <c-input-confirm
                show-icon
                @confirmed="options.item.dropdown.items.splice(dropIndex, 1)"
              />
            </td>
          </tr>
        </table>
      </div>

      <div class="mb-4 mb-3 px-3">
        <b-button
          variant="primary"
          class="text-decoration-none"
          @click="options.item.dropdown.items.push({ text: '', url: '', target: 'sameTab', delimiter: false })"
        >
          <font-awesome-icon
            :icon="['fas', 'plus']"
            size="sm"
            class="mr-1"
          />
          {{ $t("navigation.addDropdown") }}
        </b-button>
      </div>
    </td>
  </tr>
</template>

<script>
import base from './base'

export default {
  extends: base,

  data () {
    return {
      aligns: [
        { value: 'right', text: this.$t('navigation.right') },
        { value: 'left', text: this.$t('navigation.left') },
        { value: 'bottom', text: this.$t('navigation.bottom') },
        { value: 'top', text: this.$t('navigation.top') },
      ],
      targetOptions: [
        { value: 'sameTab', text: this.$t('navigation.sameTab') },
        { value: 'newTab', text: this.$t('navigation.newTab') },
      ],
    }
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    setDefaultValues () {
      this.aligns = []
      this.targetOptions = []
    },
  },
}
</script>

<style lang="scss" scoped>
th,
td {
  padding-left: 15px;
  padding-right: 15px;
}
</style>
