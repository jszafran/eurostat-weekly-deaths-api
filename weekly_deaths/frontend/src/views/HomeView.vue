<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/grave.png">
    <WelcomeScreen></WelcomeScreen>
  </div>
</template>

<script>
// @ is an alias to /src
import WelcomeScreen from '@/components/WelcomeScreen.vue'

export default {
  name: 'HomeView',
  components: {
    WelcomeScreen
  }
}
</script>

<style>
img {
  max-width: 150px;
}
</style>