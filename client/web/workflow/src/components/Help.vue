<template>
  <div>
    <b-form-group
      label-class="sticky-top bg-white pb-0 text-primary"
      class="p-2 mb-0"
    >
      <template #label>
        <h5
          class="d-flex align-items-center p-2 mb-0"
        >
          {{ $t('help:basic-controls-label') }}
          <font-awesome-icon
            :icon="['fas', 'mouse']"
            class="ml-2"
          />
        </h5>
      </template>
      <b-table
        :items="basicControls"
        fixed
        thead-class="d-none"
      />
    </b-form-group>
    <b-form-group
      label-class="sticky-top bg-white pb-0 text-primary"
      class="p-2 mb-0"
    >
      <template #label>
        <h5
          class="d-flex align-items-center pt-2 px-2 mb-0"
        >
          {{ $t('help:keyboard-shortcuts-label') }}
          <font-awesome-icon
            :icon="['fas', 'keyboard']"
            class="ml-2"
          />
        </h5>
      </template>
      <small
        class="ml-2"
      >
        {{ $t('help:mac-users') }}
      </small>
      <b-table
        :items="keyboardShortcuts"
        fixed
        thead-class="d-none"
      />
    </b-form-group>
  </div>
</template>

<script>
export default {
  data () {
    return {
      basicControls: [
        { action: this.$t('help:basic-controls.action.select'), shortcut: this.$t('help:basic-controls.shortcut.left-click.label') },
        { action: this.$t('help:basic-controls.action.select-multiple'), shortcut: this.$t('help:basic-controls.shortcut.left-click.drag') },
        { action: this.$t('help:basic-controls.action.deselect'), shortcut: this.$t('help:basic-controls.shortcut.left-click.background') },
        { action: this.$t('help:basic-controls.action.pan'), shortcut: this.$t('help:basic-controls.shortcut.right-click-drag') },
        { action: this.$t('help:basic-controls.action.zoom-in/out'), shortcut: this.$t('help:basic-controls.shortcut.mouse-wheel') },
      ],

      keyboardShortcuts: [
        { action: this.$t('help:keyboard-shortcuts.action.select-all'), shortcut: 'Ctrl + A' },
        { action: this.$t('help:keyboard-shortcuts.action.add-selection'), shortcut: `Ctrl + ${this.$t('help:keyboard-shortcuts.shortcut.left-mouse-button')}` },
        { action: this.$t('help:keyboard-shortcuts.action.delete-selected-elements'), shortcut: 'Delete or Backspace' },
        { action: this.$t('help:keyboard-shortcuts.action.reset-view-to-start-point'), shortcut: 'Ctrl + Space' },
        { action: this.$t('help:keyboard-shortcuts.action.undo'), shortcut: 'Ctrl + Z' },
        { action: this.$t('help:keyboard-shortcuts.action.redo'), shortcut: 'Ctrl + Shift + Z' },
        { action: this.$t('help:keyboard-shortcuts.action.cut'), shortcut: 'Ctrl + X' },
        { action: this.$t('help:keyboard-shortcuts.action.copy'), shortcut: 'Ctrl + C' },
        { action: this.$t('help:keyboard-shortcuts.action.paste'), shortcut: 'Ctrl + V' },
        { action: this.$t('help:keyboard-shortcuts.action.save'), shortcut: 'Ctrl + S' },
        { action: this.$t('help:keyboard-shortcuts.action.nudge'), shortcut: `${this.$t('help:keyboard-shortcuts.shortcut.arrow-keys')} Shift ${this.$t('help:keyboard-shortcuts.shortcut.adjust-distance')}` },
        { action: this.$t('help:keyboard-shortcuts.action.show-guides'), shortcut: `${this.$t('help:keyboard-shortcuts.shortcut.hold')} Alt` },
        { action: this.$t('help:keyboard-shortcuts.action.show-help'), shortcut: 'Shift + ?' },
      ],
    }
  },
}
</script>

<style>

</style>
