<script setup>
import Authenticator from './components/Authenticator.vue'
</script>

<template>
  <div @contextmenu.prevent>
    <Authenticator/>
  </div>
</template>

<style>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>
