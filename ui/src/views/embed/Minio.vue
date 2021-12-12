<template>
  <object data="http://localhost:9000/minio">
    <embed src="http://localhost:9000/minio" />
    Error: Embedded data could not be displayed.
  </object>
</template>

<script>
export default {
  name: "Minio",
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
