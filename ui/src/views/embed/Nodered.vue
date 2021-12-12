<template>
  <object data="http://localhost:1880">
    <embed src="http://localhost:1880" />
    Error: Embedded data could not be displayed.
  </object>
</template>

<script>
export default {
  name: "Nodered",
};
</script>

<style scoped>
object,
embed {
  padding-left: 56px;
  height: 100%;
  width: 100%;
}
</style>
