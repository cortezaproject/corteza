<template>
  <header>
    <div class="wrapper">
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/log-viewer">Logs</RouterLink>
      </nav>
    </div>
  </header>
  <div
    v-if="capturedError"
    class="error"
  >
    <h2>Error</h2>
    <p v-html="capturedError"></p>
  </div>
  <main
    v-else
  >
    <router-view />
  </main>
</template>

<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { onBeforeMount, ref } from 'vue'

const capturedError = ref<string | undefined>(undefined)

onBeforeMount(() => {
  //
})

</script>

<style lang="scss">
@import '@/assets/base.css';

nav {
  background: black;
  padding: 15px;
  font-size: 140%;
  a {
    margin-right: 20px;
    color: white;
    text-decoration: none;
  }
}

.error {
  background: #ffa9a9;
  padding: 40px;
  height: 100vh;
}

</style>
