<template>
  <div class="lightbox">
    <header>
      <div class="left">
        <slot name="header.left" />
      </div>
      <div class="right">
        <slot name="header.right" />
        <button
          class="closer"
          @click="$emit('close')"
        >
          &times;
        </button>
      </div>
    </header>
    <main>
      <slot name="main" />
    </main>
  </div>
</template>

<style scoped>
/* Temporary style; update when messaging is updated with bootstrap */

.lightbox {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;

  background-color: #EEF1F3;
  z-index: 2000;
}

.lightbox header {
  display: flex;
  justify-content: space-between;
  background-color: #FFFFFF;
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  box-shadow: 0 0 3px #1E1E1E41;
  z-index: 1000;
  height: 50px;
}

.lightbox header .left {
  margin-left: 10px;
}
.lightbox header .right {
  margin-right: 10px;
}

.lightbox header > div {
  display: flex;
  align-items: center;
}

.lightbox main {
  position: absolute;
  top: 50px;
  left: 0;
  right: 0;
  bottom: 0;
  overflow-y: scroll;
  z-index: 1;
}

.lightbox header button {
  background-color: transparent;
  border: none;
  outline: none;
  font-size: 25px;
  font-weight: bold;
  margin: 10px;
  cursor: pointer;
  box-sizing: border-box;
}

.lightbox header button:last-of-type {
  margin-right: 0;
}

</style>
