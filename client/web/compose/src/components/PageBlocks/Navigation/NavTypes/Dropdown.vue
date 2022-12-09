<template>
  <tr>
    <td />

    <td
      colspan="4"
      class="p-0"
    >
      <div class="d-flex">
        <th>
          <b-form-group :label="$t('navigation.fieldLabel')">
            <b-form-input
              v-model="options.item.dropdown.label"
              type="text"
            />
          </b-form-group>
        </th>

        <th>
          <b-form-group
            horizontal
            :label="$t('navigation.drop')"
          >
            <b-form-radio-group
              v-model="options.item.align"
              buttons
              button-variant="outline-secondary"
              size="sm"
              :options="aligns"
            />
          </b-form-group>
        </th>
      </div>

      <div class="d-flex align-items-center mb-4 mb-3 px-3">
        <h6 class="text-primary mb-0">
          {{ $t("navigation.dropdownItems") }}
        </h6>

        <b-button
          variant="link"
          class="text-decoration-none"
          @click="options.item.dropdown.items.push({ text: '', url: '', target: 'sameTab' })"
        >
          <font-awesome-icon
            :icon="['fas', 'plus']"
            size="sm"
            class="mr-1"
          />
          {{ $t("navigation.add") }}
        </b-button>
      </div>

      <div class="px-3">
        <table
          v-if="options.item.dropdown.items.length > 0"
          class="table table-sm table-borderless table-responsive-lg"
        >
          <tr>
            <th>
              {{ $t("navigation.text") }}
            </th>
            <th>
              {{ $t("navigation.url") }}
            </th>
            <th class="text-center">
              {{ $t('navigation.openIn') }}
            </th>
            <th class="text-center">
              {{ $t("navigation.delimiter") }}
            </th>
          </tr>

          <tr
            v-for="(item, dropIndex) in options.item.dropdown.items"
            :key="`drop-${dropIndex}`"
          >
            <th>
              <b-form-group class="mb-0">
                <b-form-input
                  v-model="item.label"
                  type="text"
                />
              </b-form-group>
            </th>

            <th>
              <b-form-group class="mb-0">
                <b-form-input
                  v-model="item.url"
                  type="text"
                />
              </b-form-group>
            </th>

            <th class="align-middle text-center">
              <b-form-group class="mb-0">
                <b-form-select
                  v-model="item.target"
                  :options="targetOptions"
                />
              </b-form-group>
            </th>

            <th class="align-middle text-center">
              <b-form-group class="mb-0">
                <b-form-checkbox
                  v-model="item.delimiter"
                  switch
                  size="sm"
                />
              </b-form-group>
            </th>

            <th class="align-middle text-center">
              <c-input-confirm
                @confirmed="options.item.dropdown.items.splice(dropIndex, 1)"
              />
            </th>
          </tr>
        </table>
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
}
</script>
