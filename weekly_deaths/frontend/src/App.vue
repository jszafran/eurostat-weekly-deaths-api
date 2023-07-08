<template>
  <div id="app">
    <nav>
<!--      <router-link to="/">Home</router-link>-->
    </nav>
    <router-view/>
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css?family=Caprasimo');
@import url('https://fonts.googleapis.com/css?family=Cardo');


#app {
  font-family: Cardo, Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #cdd5de;
}

nav {
  padding: 30px;
}

nav a {
  font-weight: bold;
  color: #2c3e50;
}

nav a.router-link-exact-active {
  color: #cdd5de;
}

body {
  background-color: #181e25;
}
</style>
<script setup>
</script>

