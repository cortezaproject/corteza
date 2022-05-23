<template>
  <section>
    <h1>Log viewer</h1>
    <table>
      <thead>
        <tr>
          <td>one</td>
          <td>level</td>
          <td>logger</td>
          <td>message</td>
        </tr>
      </thead>
      <tfoot v-if="lastRefresh">
        <tr>
          <td
            colspan="4"
            class="last-refresh"
          >
            {{ lastRefresh.toISOString().substring(11) }}
          </td>
        </tr>
      </tfoot>
      <tbody>
        <tr
          v-for="e in events"
          :key="e.index"
          :class="[`level-${e.level}`]"
        >
          <td class="ts">{{ e.ts.toISOString().substring(11) }}</td>
          <td class="level">{{ e.level }}</td>
          <td class="logger">{{ e.logger }}</td>
          <td class="msg">
            {{ e.msg }}
            <pre v-if="e.extra">{{ e.extra }}</pre>
          </td>
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
const lastRefresh = ref<Date | undefined>()


onMounted(() => {
  events.value = []
  fetch().then((interval: boolean) => {
    let ih: number
    if (!interval) {
      return
    }

    ih = setInterval(async () => {
      let after: number | undefined

      if (events.value.length > 0) {
        after = events.value[events.value.length - 1].index
      }

      const ok = await fetch(after)
      if (!ok) {
        clearInterval(ih)
      }
    }, 2000)
  }).catch((err ) => {
    alert(err)
  })

})

async function fetch (after?: number): Promise<boolean> {
  return fetchLoggedEvents({ after, limit: -1 })
    .then(ee => {
      events.value.push(...ee)
      lastRefresh.value = new Date()
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
        border-top: 1px dotted silver;
        vertical-align: top;

        &.msg pre {
          color: gray;
        }
      }
    }
  }

  tfoot {
    .last-refresh {
      padding: 20px;
      text-align: center;
    }
  }
}
</style>
