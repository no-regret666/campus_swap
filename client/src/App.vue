<script setup>
import Navbar from './components/Navbar.vue'
import Toast from './components/Toast.vue'
</script>

<template>
  <div id="app-root">
    <Navbar />
    <main class="main-content">
      <router-view />
    </main>
    <Toast />
  </div>
</template>

<style scoped>
.main-content {
  padding-top: 64px;
  min-height: calc(100vh - 64px);
}
</style>
