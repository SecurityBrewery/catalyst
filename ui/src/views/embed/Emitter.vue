<template>
  <object data="http://localhost:9001/keygen">
    <embed src="http://localhost:9001/keygen" />
    Error: Embedded data could not be displayed.
  </object>
</template>

<script>
export default {
  name: "Emitter",
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
