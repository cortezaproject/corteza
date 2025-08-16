<template>
  <div
    v-if="items.length"
    class="mention-dropdown"
  >
    <button
      v-for="(item, index) in items"
      :key="index"
      type="button"
      :class="[
        'mention-option',
        { 'mention-option--highlighted': index === selectedIndex }
      ]"
      @click="handleClick(index, $event)"
      @mouseenter="selectedIndex = index"
    >
      {{ getDisplayName(item) }}
    </button>
  </div>
</template>

<script>
export default {
  name: 'MentionList',

  props: {
    items: {
      type: Array,
      required: true,
    },

    command: {
      type: Function,
      required: true,
    },
  },

  data() {
    return {
      selectedIndex: 0,
    }
  },

  watch: {
    items() {
      this.selectedIndex = 0
    },
  },

  methods: {
    onKeyDown({ event }) {
      if (event.key === 'ArrowUp') {
        this.upHandler()
        return true
      }

      if (event.key === 'ArrowDown') {
        this.downHandler()
        return true
      }

      if (event.key === 'Enter') {
        this.enterHandler()
        return true
      }

      return false
    },

    upHandler() {
      this.selectedIndex = ((this.selectedIndex + this.items.length) - 1) % this.items.length
    },

    downHandler() {
      this.selectedIndex = (this.selectedIndex + 1) % this.items.length
    },

    enterHandler() {
      this.selectItem(this.selectedIndex)
    },

    handleClick(index, event) {
      this.selectItem(index)
    },

    selectItem(index) {
      const item = this.items[index]

      if (item) {
        this.command({
          id: item.userID,
          label: this.getDisplayName(item),
        })
      }
    },

    getDisplayName(user) {
      const { name, username, email, userID } = user
      return name || username || email || userID
    },
  },
}
</script>

<style lang="scss" scoped>
.mention-dropdown {
  background: var(--white);
  border: 1px solid var(--extra-light);
  border-radius: 0.25rem;
  box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
  max-height: 200px;
  overflow-y: auto;
  font-size: 0.9rem;
  font-family: var(--font-regular);
}

// Mirroring CInputSelect's vs__dropdown-option styles
.mention-option {
  display: block;
  width: 100%;
  background: var(--white);
  color: var(--black);
  padding: 0.5rem 0.75rem;
  border: none;
  text-align: left;
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;

  &:hover,
  &.mention-option--highlighted {
    background: var(--light);
    color: var(--black);
  }

  &:active,
  &:focus {
    color: var(--white);
    background-color: var(--primary);
    outline: none;
  }
}
</style>
