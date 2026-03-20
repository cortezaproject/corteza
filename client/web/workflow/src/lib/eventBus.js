/**
 * Simple event bus to replace Vue 2's $root.$emit / $root.$on pattern.
 * Used for cross-component communication in the workflow editor.
 */
import mitt from 'mitt'

const eventBus = mitt()

export default eventBus
