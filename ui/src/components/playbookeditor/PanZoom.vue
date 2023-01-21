<template>
  <div id="graphwrapper">
    <div id="gab">
      <slot name="actionbar"/>
    </div>
    <v-toolbar
        id="gtb"
        class="ma-2"
        floating
        dense
    >
      <v-btn @click.prevent.stop="reset" title="Reset" icon>
        <v-icon color="#000">mdi-image-filter-center-focus</v-icon>
      </v-btn>
      <v-btn @click.prevent.stop="zoomIn" title="Zoom in" :disabled="isMaxZoom" icon>
        <v-icon color="#000">mdi-magnify-plus-outline</v-icon>
      </v-btn>
      <v-btn @click.prevent.stop="zoomOut" title="Zoom out" :disabled="isMinZoom" icon>
        <v-icon color="#000">mdi-magnify-minus-outline</v-icon>
      </v-btn>
      <slot name="toolbar"/>
    </v-toolbar>
    <div id="panzoom">
      <slot/>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref, defineProps, defineExpose} from "vue";
import createPanZoom from "panzoom";
import * as panzoom from "panzoom";

const props = defineProps<{
  config: panzoom.PanZoomOptions
}>();

const panZoom = ref<panzoom.PanZoom | null>(null);

const zoomLevel = ref<number>(1);

const minZoom = ref<number>(0.5);
const maxZoom = ref<number>(1.5);

const isMaxZoom = computed(() => {
  if (zoomLevel.value) {
    return zoomLevel.value >= maxZoom.value;
  }
  return false;
});

const isMinZoom = computed(() => {
  if (zoomLevel.value) {
    return zoomLevel.value <= minZoom.value;
  }
  return false;
});

const initialZoom = ref<number>(1);

onMounted(() => {
  const canvas = document.getElementById("panzoom")
  if (!canvas) {
    throw new Error("No element with id panzoom")
  }

  const firstChild = canvas.firstElementChild
  if (!firstChild) {
    throw new Error("No child element")
  }

  const startX = canvas.getBoundingClientRect().width / 2 - firstChild.getBoundingClientRect().width / 2
  const startY = canvas.getBoundingClientRect().height / 2 - firstChild.getBoundingClientRect().height / 2

  initialZoom.value = props.config.initialZoom ? props.config.initialZoom : 1
  minZoom.value = props.config.initialZoom / 2
  maxZoom.value = props.config.initialZoom * 2

  panZoom.value = createPanZoom(canvas, {
    ...props.config,
    zoomDoubleClickSpeed: 1, // disable double click zoom
    autocenter: true,
    initialX: startX,
    initialY: startY,
    minZoom: minZoom.value,
    maxZoom: maxZoom.value,
  });

  panZoom.value.on("zoom", (e: panzoom.PanZoom) => {
    zoomLevel.value = e.getTransform().scale;
  });

  reset(false);
});

const ZOOM_FACTOR = 0.255

const zoomIn = () => {
  if (!panZoom.value) return
  const currentZoom = panZoom.value.getTransform().scale
  panZoom.value.smoothZoomAbs(0, 0, currentZoom + ZOOM_FACTOR)
}

const zoomOut = () => {
  if (!panZoom.value) return
  const currentZoom = panZoom.value.getTransform().scale
  panZoom.value.smoothZoomAbs(0, 0, currentZoom - ZOOM_FACTOR)
}

const reset = (smooth = true) => {
  if (!panZoom.value) return

  const canvas = document.getElementById("panzoom")
  if (!canvas) {
    throw new Error("No element with id panzoom")
  }

  const firstChild = canvas.firstElementChild
  if (!firstChild) {
    throw new Error("No child element")
  }

  const startX = canvas.getBoundingClientRect().width / 2 - firstChild.getBoundingClientRect().width / 2
  const startY = canvas.getBoundingClientRect().height / 2 - firstChild.getBoundingClientRect().height / 2

  panZoom.value.pause()
  if (smooth) {
    panZoom.value.smoothZoomAbs(0, 0, initialZoom.value)
    panZoom.value.smoothMoveTo(startX * 0.5, startY)
  } else {
    panZoom.value.zoomTo(0, 0, initialZoom.value)
    panZoom.value.moveTo(startX * 0.5, startY)
  }
  panZoom.value.resume()
}

defineExpose({
  reset,
  zoomIn,
  zoomOut,
  zoomLevel,
  isMaxZoom,
  isMinZoom,
});
</script>

<style>
#graphwrapper {
  position: relative;
  width: 100%;
  height: 100%;
}

#gtb .v-toolbar__content {
  padding: 0;
}

#panzoom {
  width: 100%;
  height: 100%;
  position: relative;
  cursor: move;
}
</style>