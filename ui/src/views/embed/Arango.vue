<template>
  <object data="http://localhost:8003">
    <embed src="http://localhost:8003" />
    Error: Embedded data could not be displayed.
  </object>
</template>

<script>
export default {
  name: "Arango",
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
