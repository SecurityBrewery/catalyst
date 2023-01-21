<template>
  <svg
      class="pe-graph"
      :width="width * scale"
      :height="height * scale"
      :viewBox="`${minX - config.graphPadding} ${minY - config.graphPadding} ${maxX - minX + config.boxWidth + config.graphPadding * 2} ${maxY - minY + config.boxHeight + config.graphPadding * 2}`"
  >
    <defs>
      <path
          id="cog"
          d="M12,15.5A3.5,3.5 0 0,1 8.5,12A3.5,3.5 0 0,1 12,8.5A3.5,3.5 0 0,1 15.5,12A3.5,3.5 0 0,1 12,15.5M19.43,12.97C19.47,12.65 19.5,12.33 19.5,12C19.5,11.67 19.47,11.34 19.43,11L21.54,9.37C21.73,9.22 21.78,8.95 21.66,8.73L19.66,5.27C19.54,5.05 19.27,4.96 19.05,5.05L16.56,6.05C16.04,5.66 15.5,5.32 14.87,5.07L14.5,2.42C14.46,2.18 14.25,2 14,2H10C9.75,2 9.54,2.18 9.5,2.42L9.13,5.07C8.5,5.32 7.96,5.66 7.44,6.05L4.95,5.05C4.73,4.96 4.46,5.05 4.34,5.27L2.34,8.73C2.21,8.95 2.27,9.22 2.46,9.37L4.57,11C4.53,11.34 4.5,11.67 4.5,12C4.5,12.33 4.53,12.65 4.57,12.97L2.46,14.63C2.27,14.78 2.21,15.05 2.34,15.27L4.34,18.73C4.46,18.95 4.73,19.03 4.95,18.95L7.44,17.94C7.96,18.34 8.5,18.68 9.13,18.93L9.5,21.58C9.54,21.82 9.75,22 10,22H14C14.25,22 14.46,21.82 14.5,21.58L14.87,18.93C15.5,18.67 16.04,18.34 16.56,17.94L19.05,18.95C19.27,19.03 19.54,18.95 19.66,18.73L21.66,15.27C21.78,15.05 21.73,14.78 21.54,14.63L19.43,12.97Z"/>
      <path
          id="clipboard-outline"
          d="M19,3H14.82C14.4,1.84 13.3,1 12,1C10.7,1 9.6,1.84 9.18,3H5A2,2 0 0,0 3,5V19A2,2 0 0,0 5,21H19A2,2 0 0,0 21,19V5A2,2 0 0,0 19,3M12,3A1,1 0 0,1 13,4A1,1 0 0,1 12,5A1,1 0 0,1 11,4A1,1 0 0,1 12,3M7,7H17V5H19V19H5V5H7V7Z"/>
      <path
          id="keyboard"
          d="M19,10H17V8H19M19,13H17V11H19M16,10H14V8H16M16,13H14V11H16M16,17H8V15H16M7,10H5V8H7M7,13H5V11H7M8,11H10V13H8M8,8H10V10H8M11,11H13V13H11M11,8H13V10H11M20,5H4C2.89,5 2,5.89 2,7V17A2,2 0 0,0 4,19H20A2,2 0 0,0 22,17V7C22,5.89 21.1,5 20,5Z"/>
      <path
          id="star"
          d="M12,17.27L18.18,21L16.54,13.97L22,9.24L14.81,8.62L12,2L9.19,8.62L2,9.24L7.45,13.97L5.82,21L12,17.27Z"/>
    </defs>
    <g class="pe-links">
      <path v-for="(link, index) in links" :key="index" :d="link.path"/>
    </g>
    <g class="pe-activelinks">
      <path
          v-for="(link, index) in links"
          :key="index"
          :d="link.path"
          :class="{
            'hovered': !selectedNode &&  (link.source === hoverNode || link.target === hoverNode),
            'selected': link.source === selectedNode || link.target === selectedNode
          }"
      />
    </g>
    <g class="pe-nodes">
      <g
          v-for="(node, index) in positionedNodes"
          :key="index">
        <g
            v-if="node.type === 'start'"
            :transform="`translate(${props.horizontal ? node.x + config.boxWidth : node.x + (config.boxWidth / 2)}, ${node.y})`"
            class="start"
        >
          <circle
              :transform="`translate(0, ${config.boxHeight / 2})`"
              r="20"
          />
          <g
              class="icon"
              :transform="`translate(${0 - 12}, ${(config.boxHeight - 24) / 2})`"
          >
            <use href="#star" fill="white"/>
          </g>
        </g>
        <g
            v-else
            :transform="`translate(${node.x}, ${node.y})`"
            :class="{
            'pe-node': true,
            'start': node.type === 'start',
            'hovered': node.id === hoverNode,
            'selected': node.id === selectedNode,
            'unhovered': hoverNode && !selectedNode && node.id !== hoverNode,
            'unselected': selectedNode && node.id !== selectedNode,
          }"
            @mouseover="hoverNode = node.id"
            @mouseout="hoverNode = null"
            @click="selectedNode === node.id ? selectedNode = null : selectedNode = node.id"
        >
          <rect
              class="pe-box"
              :width="config.boxWidth"
              :height="config.boxHeight"
              :rx="config.boxRadius"
              :ry="config.boxRadius"
          />
          <g
              class="pe-icon"
              :transform="`translate(10, ${(config.boxHeight - 24) / 2})`"
          >
            <use v-if="node.type === 'automation'" href="#cog"/>
            <use v-else-if="node.type === 'start'" href="#star"/>
            <use v-else-if="node.type === 'input'" href="#keyboard"/>
            <use v-else href="#clipboard-outline"/>
          </g>
          <text
              class="pe-text"
              :x="config.boxWidth / 2"
              :y="config.boxHeight / 2"
          >
            {{ node.label ? node.label : node.id }}
          </text>
          <g
              class="add"
              :transform="`translate(${config.boxWidth - 34}, ${(config.boxHeight - 24) / 2})`"
          >
            <use href="#star" fill="none"/>
          </g>
        </g>
      </g>
    </g>
    <g class="pe-connectors">
      <g
          v-for="(link, index) in links"
          :key="index"
          :class="{
            'pe-connector': true,
            'hovered': !selectedNode && (link.source === hoverNode || link.target === hoverNode),
            'selected': link.source === selectedNode || link.target === selectedNode
          }"
      >
        <circle
            v-if="link.source !== 'start'"
            :cx="link.start.x"
            :cy="link.start.y"
            :r="4"/>
        <circle
            :cx="link.end.x"
            :cy="link.end.y"
            :r="7"/>
        <circle
            :cx="link.end.x"
            :cy="link.end.y"
            :r="5"
            fill="white"/>
      </g>
    </g>
  </svg>
</template>

<script setup lang="ts">
import * as d3 from "d3";
import {digl} from "@crinkles/digl";
import {Edge as DiglEdge, Node as DiglNode, Position, Rank} from "@crinkles/digl/dist/types";
import {computed, defineEmits, ref, Ref, defineProps, ComputedRef} from "vue";

interface Config {
  graphPadding: number;
  boxWidth: number;
  boxHeight: number;
  boxMarginX: number;
  boxMarginY: number;
  boxRadius: number;
  lineDistance: number;
}

const props = defineProps({
  playbook: {
    type: Object,
    required: true
  },
  selected: {
    type: String,
    required: false,
    default: null
  },
  horizontal: {
    type: Boolean,
    required: false,
    default: false
  },
  scale: {
    type: Number,
    required: false,
    default: 1
  }
})

const config = ref({
  graphPadding: 10,
  boxWidth: 220,
  boxHeight: 40,
  boxMarginX: 60,
  boxMarginY: 70,
  boxRadius: 20,
  lineDistance: 0,
});

interface Edge extends DiglEdge {
  label?: { text: string, x: number, y: number };
}

interface Node extends DiglNode {
  type: 'automation' | 'input' | 'task' | 'start';
}

const edges = computed(() => {
  const edges: Array<Edge> = [];

  for (const key in props.playbook.tasks) {
    for (const next in props.playbook.tasks[key].next) {
      edges.push({
        source: key,
        target: next,
        label: props.playbook.tasks[key].next[next]
      });
    }
  }

  const rootNodes = nodes.value.filter(node => edges.every(edge => edge.target !== node.id));
  for (const node of rootNodes) {
    if (node.id !== 'start') {
      edges.push({
        source: 'start',
        target: node.id
      });
    }
  }

  return edges;
});

const nodes: ComputedRef<Array<Node>> = computed(() => {
  const nodes = [{
    id: "start",
    label: "Start",
    type: "start" as 'automation' | 'input' | 'task' | 'start',
  }];
  for (const key in props.playbook.tasks) {
    nodes.push({
      id: key,
      label: props.playbook.tasks[key].name,
      type: props.playbook.tasks[key].type
    });
  }
  return nodes;
});


const hoverNode: Ref<string | null> = ref(null);

const emits = defineEmits(['update:selected']);
const selectedNode = computed({
  get: () => props.selected,
  set: (newVal) => {
    emits('update:selected', newVal);
  }
});

const rankses: Ref<Array<Array<Rank>>> = computed(() => {
  return digl(edges.value);
});

const ranks = computed(() => {
  return rankses.value.length > 0 ? rankses.value[0] : [];
});

const positionedNodes = computed(() => {
  return positioning(config.value, props.horizontal, nodes.value, ranks.value);
});

const width = computed(() => {
  return maxX.value - minX.value + config.value.graphPadding * 2 + config.value.boxWidth;
});

const height = computed(() => {
  return maxY.value - minY.value + config.value.graphPadding * 2 + config.value.boxHeight;
});

const minX = computed(() => {
  if (props.horizontal) {
    return Math.min(...positionedNodes.value.map(node => node.x)) + config.value.boxWidth - 22;
  }
  return Math.min(...positionedNodes.value.map(node => node.x));
});

const minY = computed(() => {
  if (!props.horizontal) {
    return Math.min(...positionedNodes.value.map(node => node.y)) + config.value.boxHeight - 30;
  }
  return Math.min(...positionedNodes.value.map(node => node.y));
});

const maxX = computed(() => {
  return Math.max(...positionedNodes.value.map(node => node.x));
});

const maxY = computed(() => {
  return Math.max(...positionedNodes.value.map(node => node.y)) + 50;
});

const links = computed(() => {
  return edges.value.map(edge => {
    const source = positionedNodes.value.find(node => node.id === edge.source);
    const target = positionedNodes.value.find(node => node.id === edge.target);

    if (!source || !target) return;

    // index within rank
    const sourceIndex = ranks.value.find(rank => rank.includes(edge.source))?.indexOf(edge.source);
    if (sourceIndex === undefined) {
      throw new Error(`sourceIndex is undefined for ${edge.source}`);
    }

    const path = props.horizontal ?
        horizontalConnectionLine(source, target, sourceIndex, config.value) :
        verticalConnectionLine(source, target, sourceIndex, config.value);

    const start = props.horizontal ?
        {x: source.x + config.value.boxWidth, y: source.y + config.value.boxHeight / 2} :
        {x: source.x + config.value.boxWidth / 2, y: source.y + config.value.boxHeight};

    const end = props.horizontal ?
        {x: target.x, y: target.y + config.value.boxHeight / 2} :
        {x: target.x + config.value.boxWidth / 2, y: target.y};

    return {
      source: edge.source,
      target: edge.target,
      path: path.toString(),
      start: start,
      end: end,
    };
  }).filter(link => link);
});

interface PositionedNode extends Node, Position {
}

function positioning(
    config: Config,
    horizontal: boolean,
    nodes: Node[],
    ranks: Rank[]
): PositionedNode[] {
  const _nodes: PositionedNode[] = [];
  const _h = horizontal;

  ranks.forEach((rank, i) => {
    const xStart = _h
        ? (config.boxWidth + config.boxMarginX) * i
        : -0.5 * (rank.length - 1) * (config.boxWidth + config.boxMarginX);
    const yStart = _h
        ? -0.5 * (rank.length - 1) * (config.boxHeight + config.boxMarginY)
        : (config.boxHeight + config.boxMarginY) * i;

    rank.forEach((nodeId, nIndex) => {
      const _node: Node = nodes.find((n) => n.id == nodeId) as Node;
      if (!_node) return;
      const x = _h ? xStart : xStart + (config.boxWidth + config.boxMarginX) * nIndex;
      const y = _h ? yStart + (config.boxHeight + config.boxMarginY) * nIndex : yStart;
      _nodes.push({..._node, x, y});
    });
  });

  return _nodes;
}

function verticalConnectionLine(source: PositionedNode, target: PositionedNode, sourceIndex: number, config: Config) {
  const sourceBottomCenter = {
    x: source.x + config.boxWidth / 2,
    y: source.y + config.boxHeight,
  };

  const targetTopCenter = {
    x: target.x + config.boxWidth / 2,
    y: target.y,
  };

  const path = d3.path();
  path.moveTo(sourceBottomCenter.x, sourceBottomCenter.y);

  const lineCurve = config.boxMarginY / 2;

  if (sourceBottomCenter.x == targetTopCenter.x) {
    path.lineTo(targetTopCenter.x, targetTopCenter.y);
  } else if (sourceBottomCenter.x < targetTopCenter.x) {
    if (target.y !== source.y + config.boxHeight + config.boxMarginY) {
      path.lineTo(sourceBottomCenter.x, target.y - config.boxMarginY);
    }
    sourceBottomCenter.y = target.y - config.boxMarginY;
    path.quadraticCurveTo(
        sourceBottomCenter.x, sourceBottomCenter.y + lineCurve + sourceIndex * config.lineDistance,
        sourceBottomCenter.x + lineCurve, sourceBottomCenter.y + lineCurve + sourceIndex * config.lineDistance,
    );
    path.lineTo(targetTopCenter.x - lineCurve, targetTopCenter.y - lineCurve + sourceIndex * config.lineDistance);
    path.quadraticCurveTo(
        targetTopCenter.x, targetTopCenter.y - lineCurve + sourceIndex * config.lineDistance,
        targetTopCenter.x, targetTopCenter.y,
    );
  } else {
    if (target.y !== source.y + config.boxHeight + config.boxMarginY) {
      path.lineTo(sourceBottomCenter.x, target.y - config.boxMarginY);
    }
    sourceBottomCenter.y = target.y - config.boxMarginY;
    path.quadraticCurveTo(
        sourceBottomCenter.x, sourceBottomCenter.y + lineCurve + sourceIndex * config.lineDistance,
        sourceBottomCenter.x - lineCurve, sourceBottomCenter.y + lineCurve + sourceIndex * config.lineDistance,
    );
    path.lineTo(targetTopCenter.x + lineCurve, targetTopCenter.y - lineCurve + sourceIndex * config.lineDistance);
    path.quadraticCurveTo(
        targetTopCenter.x, targetTopCenter.y - lineCurve + sourceIndex * config.lineDistance,
        targetTopCenter.x, targetTopCenter.y,
    );
  }
  return path;
}

function horizontalConnectionLine(source: PositionedNode, target: PositionedNode, sourceIndex: number, config: Config) {
  const sourceRightCenter = {
    x: source.x + config.boxWidth,
    y: source.y + config.boxHeight / 2,
  };

  const targetLeftCenter = {
    x: target.x,
    y: target.y + config.boxHeight / 2,
  };

  const path = d3.path();
  path.moveTo(sourceRightCenter.x, sourceRightCenter.y);

  const lineCurve = config.boxMarginX / 2;

  if (sourceRightCenter.y == targetLeftCenter.y) {
    path.lineTo(targetLeftCenter.x, targetLeftCenter.y);
  } else if (sourceRightCenter.y < targetLeftCenter.y) {
    if (target.x !== source.x + config.boxWidth + config.boxMarginX) {
      path.lineTo(target.x - config.boxMarginX, sourceRightCenter.y);
    }
    sourceRightCenter.x = target.x - config.boxMarginX;
    path.quadraticCurveTo(
        sourceRightCenter.x + lineCurve + sourceIndex * config.lineDistance, sourceRightCenter.y,
        sourceRightCenter.x + lineCurve + sourceIndex * config.lineDistance, sourceRightCenter.y + lineCurve,
    );
    path.lineTo(targetLeftCenter.x - lineCurve + sourceIndex * config.lineDistance, targetLeftCenter.y - lineCurve);
    path.quadraticCurveTo(
        targetLeftCenter.x - lineCurve + sourceIndex * config.lineDistance, targetLeftCenter.y,
        targetLeftCenter.x, targetLeftCenter.y,
    );
  } else {
    if (target.x !== source.x + config.boxWidth + config.boxMarginX) {
      path.lineTo(target.x - config.boxMarginX, sourceRightCenter.y);
    }
    sourceRightCenter.x = target.x - config.boxMarginX;
    path.quadraticCurveTo(
        sourceRightCenter.x + lineCurve + sourceIndex * config.lineDistance, sourceRightCenter.y,
        sourceRightCenter.x + lineCurve + sourceIndex * config.lineDistance, sourceRightCenter.y - lineCurve,
    );
    path.lineTo(targetLeftCenter.x - lineCurve + sourceIndex * config.lineDistance, targetLeftCenter.y + lineCurve);
    path.quadraticCurveTo(
        targetLeftCenter.x - lineCurve + sourceIndex * config.lineDistance, targetLeftCenter.y,
        targetLeftCenter.x, targetLeftCenter.y,
    );
  }
  return path;
}
</script>

<style>
svg .pe-node.unhovered {
  opacity: 0.5;
}

svg .pe-node .pe-box {
  fill: white;
  stroke: #4C566A;
  stroke-width: 1px;
  filter: drop-shadow(0 1.5px 1.5px rgba(0, 0, 0, .15));
}

svg .pe-node.selected .pe-box {
  stroke: #88C0D0;
  stroke-width: 2;
}

svg .pe-node text {
  text-anchor: middle;
  dominant-baseline: middle;
  font-size: 16px;
}

svg .pe-node:hover .pe-box {
  filter: drop-shadow(0 0 0.1rem #333);
}

svg .pe-node:hover .pe-box,
svg .pe-node:hover text {
  cursor: pointer;
}

svg .pe-node:hover .add use {
  /* fill: red !important; */
}

svg .pe-links path {
  stroke: black;
  fill: none;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
  filter: drop-shadow(0 2px 2px rgba(0, 0, 0, .15));
}

svg .pe-activelinks path {
  stroke: none;
  fill: none;
  stroke-width: 2;
  stroke-linejoin: round;
}

svg .pe-activelinks path.selected,
svg .pe-activelinks path.hovered {
  stroke-width: 3;
  stroke-linejoin: round;
  stroke: #88C0D0;
  filter: drop-shadow(0 0 0.1rem #88C0D0);
}

</style>
