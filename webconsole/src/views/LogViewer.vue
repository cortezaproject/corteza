<template>
  <section>
    <h1>Log viewer</h1>
    <table>
      <thead>
        <tr>
          <td>one</td>
          <td>level</td>
          <td>message</td>
          <td>logger</td>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="e in events"
          :key="e.index"
          :class="[`level-${e.level}`]"
        >
          <td class="ts">{{ e.ts.toISOString().substring(11) }}</td>
          <td class="level">{{ e.level }}</td>
          <td class="msg">{{ e.msg }}</td>
          <td class="logger">{{ e.logger }}</td>
        </tr>
      </tbody>
    </table>
  </section>
</template>
<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { fetchLoggedEvents  } from '@/libs/logs'
import type { LogEntry } from '@/libs/logs'

let intervalHandler: number | undefined
const events = ref<Array<LogEntry>>([])


onMounted(() => {
  fetch().then((interval: boolean) => {
    if (!interval) {
      return
    }

    setInterval(() => {
      let after: number | undefined

      if (events.value.length > 0) {
        after = events.value[events.value.length - 1].index
      }

      fetchLoggedEvents({ after, limit: -1 })
        .then(ee => events.value.push(...ee))

    }, 2000)
  })

})

async function fetch (): Promise<boolean> {
  return fetchLoggedEvents({ limit: -1 })
    .then(ee => {
      events.value = ee
      return true
    })
    .catch((err) => {
      alert(err)
      return false
    })
}

onUnmounted(() => {
  if (intervalHandler) {
    clearInterval(intervalHandler)
  }
})


</script>
<style lang="scss" scoped>
table {
  width: 100vw;
  font-family: monospace;

  tbody {
    tr {
      &.level-warn {
        background: gold;
      }

      &.level-debug .msg {
        color: gray;
      }

      td {
        padding: 3px;
        margin: 0;
        border: 0;
      }
    }

  }
}
</style>
