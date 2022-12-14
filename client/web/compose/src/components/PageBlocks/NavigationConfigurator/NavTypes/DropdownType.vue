<template>
  <div>
    <b-row class="mt-5">
      <b-col
        cols="6"
        sm="6"
      >
        <b-form-group
          horizontal
          :label="$t('navigation.drop')"
        >
          <b-form-radio-group
            v-model="column.options.dropdownType"
            buttons
            button-variant="outline-secondary"
            size="sm"
            name="drop"
            :options="dropTypes"
          />
        </b-form-group>
      </b-col>
      <b-col
        cols="6"
        sm="6"
      >
        <b-form-group label="Text">
          <b-form-input
            v-model="column.options.dropdownText"
            type="text"
            size="sm"
          />
        </b-form-group>
      </b-col>
    </b-row>
    <b-row
      v-for="(item, dropIndex) in column.options.dropdown"
      :key="`drop-${dropIndex}`"
      align-v="center"
    >
      <b-col cols="2">
        <b-form-group label="Delimiter">
          <b-form-checkbox
            v-model="item.delimiter"
            switch
          />
        </b-form-group>
      </b-col>
      <b-col cols="4">
        <b-form-group label="Text">
          <b-form-input
            v-model="item.text"
            type="text"
            size="sm"
          />
        </b-form-group>
      </b-col>
      <b-col cols="4">
        <b-form-group label="URL">
          <b-form-input
            v-model="item.url"
            type="text"
            size="sm"
          />
        </b-form-group>
      </b-col>
      <b-col
        cols="2"
        class="d-flex mt-4 justify-content-around"
      >
        <c-input-confirm
          variant="link"
          size="lg"
          button-class="text-dark px-0"
          @confirmed="column.options.dropdown.splice(dropIndex, 1)"
        />
      </b-col>
    </b-row>

    <div class="d-flex align-items-center">
      <b-button
        class="text-decoration-none"
        size="sm"
        @click="column.options.dropdown.push({ text: '', url: '' })"
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
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    column: {
      type: Object,
      required: true,
    },
  },

  data () {
    return {
      dropTypes: [
        { value: 'right', text: this.$t('navigation.right') },
        { value: 'left', text: this.$t('navigation.left') },
        { value: 'top', text: this.$t('navigation.top') },
        { value: 'bottom', text: this.$t('navigation.bottom') },
      ],
    }
  },
}
</script>
