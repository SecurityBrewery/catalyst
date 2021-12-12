<template>
  <div>
    <div v-for="(value, key) in json" :key="key">
      <dl v-if="typeof(value) == 'object'">
        <dt>{{ key }}</dt>
        <dd><JSONHTML :json="value"></JSONHTML></dd>
      </dl>
      <span v-else>
        <b>{{ key }}</b>: {{ value }}
      </span>
    </div>
  </div>
</template>

<script>
export default {
  name: "JSONHTML",
  props: ["json"]
}
</script>
