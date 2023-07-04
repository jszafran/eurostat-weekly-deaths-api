<template>
  <div class="home">

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
  },
  title: 'Eurostat Weekly Deaths App'
}
</script>

<style>
img {
  max-width: 150px;
}

.death {
  background: #3c78a8;
  max-width: 150px;
  margin: auto;
  padding: 40px;
  border-radius: 50%;
}
</style>