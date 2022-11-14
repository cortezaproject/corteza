<script>
import base from './base'

/**
 * Helper to find a label for the given value
 * @param {String} v Value in question
 * @param {Array<Object>} options Available options
 * @returns {String|undefined}
 */
function findLabel (v, options) {
  return (options.find(({ value }) => value === v) || {}).text
}

export default {
  extends: base,

  computed: {
    /**
     * Overwrite default; allow values to resolve to their labels
     * @returns {String|Array<String>}
     */
    value () {
      let v
      if (this.field.isSystem) {
        v = this.record[this.field.name]
      }
      v = this.record ? this.record.values[this.field.name] : undefined

      if (this.field.isMulti) {
        return v.map(v => findLabel(v, this.field.options.options) || v)
      } else {
        return findLabel(v, this.field.options.options) || v
      }
    },
  },
}
</script>
