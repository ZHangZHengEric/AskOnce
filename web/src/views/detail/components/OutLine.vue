<template>
  <div class="ml-20 border border-[#F5F5F5] border-solid backdrop-blur-l bg-white/30	 rounded-xl" id="right_tree">
    <div class="flex bg-[#F5F5F5] rounded-t-xl pt-4 lg:pl-3 pr-3 pb-4">
      <div>{{ $t('detail.outline') }}</div>
      <svg-icon icon-class="detail_tree" class="w-4 h-4 mt-1 ml-10 cursor-pointer" @click="showTree(false)"></svg-icon>
      <svg-icon icon-class="detail_list" class="w-4 h-4  mt-1 ml-6 cursor-pointer" @click="showTree(true)"></svg-icon>
      <div class="flex-1"></div>
      <svg-icon icon-class="close" class="w-3 h-3 capitalize mt-1.5 cursor-pointer hidden lg:block"
                @click="hideOutLine"></svg-icon>
    </div>
    <div v-if="showOutLine">
      <div class="mt-4 px-3 h-full overflow-y-scroll scrollbar-hidden"
           v-show="data.showTree"
           :style="{width:data.chartWidth+'px',height:data.chartHeight+'px'}">
        <div v-for="(h1,i1) in list" :key="i1" class="text-size16 font-[500] text-color333 mb-6">
          <div class="mb-3 text-default">{{ h1.content }}</div>
          <template v-if="h1.children">
            <div v-for="(h2,i2) in h1.children" :key="i2" class="pl-2 text-size16 text-normal">
              <div class="mb-2">• {{ h2.content }}</div>
              <template v-if="h2.children">
                <div v-for="(h3,i3) in h2.children" :key="i3" class="pl-4 text-size16 text-color666">
                  <div class="mb-2">• {{ h3.content }}</div>
                </div>
              </template>
            </div>
          </template>
        </div>
      </div>
      <div v-show="!data.showTree" id="chart"
           :style="{width:data.chartWidth+'px',height:data.chartHeight+'px'}"></div>
    </div>
    <div v-else class="text-center mt-10 text-default text-size18"
         :style="{width:data.chartWidth+'px',height:data.chartHeight+'px'}">
      大纲生成中...
    </div>
  </div>
</template>
<script setup>
import {defineProps, computed, reactive, onMounted, defineEmits, onUnmounted} from 'vue'
import G6 from '@antv/g6';

const data = reactive({
  showTree: true,
  myChart: null,
  chartWidth: 0,
  chartHeight: 0,
  container: '',
  graph: null,
})
onMounted(() => {
  data.chartWidth = document.getElementById('right_tree').clientWidth
  data.chartHeight = document.body.clientHeight - 160
  window.addEventListener('resize', resize)
})

onUnmounted(() => {
  window.removeEventListener('resize', resize)
})

const resize = () => {
  data.chartWidth = document.getElementById('right_tree').clientWidth
  data.chartHeight = document.body.clientHeight - 160
  setTimeout(() => {
    if (!data.graph || data.graph.get('destroyed')) return;
    if (!data.chartWidth || !data.chartHeight) return;
    data.graph.changeSize(data.chartWidth, data.chartHeight);
    data.graph.fitCenter()
  }, 100);
}

const emits = defineEmits([
  'hide'
])

const showTree = (show) => {
  if (!props.showOutLine) {
    return
  }
  data.showTree = show
  if (!show) {
    initChart()
  }
}

const hideOutLine = () => {
  emits('hide')
}

const initChart = () => {
  if (data.graph) {
    data.graph.destroy()
  }
  const container = document.getElementById('chart');
  const width = data.chartWidth;
  const height = data.chartHeight;
  data.graph = new G6.TreeGraph({
    container: container,
    width,
    height,
    modes: {
      default: [
        {
          type: 'collapse-expand',
          onChange: function onChange(item, collapsed) {
            const data = item.get('model');
            data.collapsed = collapsed;
            return true;
          },
        },
        'drag-canvas',
        'zoom-canvas',
      ],
    },
    defaultNode: {
      size: 8,
      anchorPoints: [
        [0, 0.5],
        [1, 0.5],
      ],
    },
    defaultEdge: {
      type: 'cubic-horizontal',
    },
    layout: {
      type: 'mindmap',
      direction: 'H',
      getHeight: () => {
        return 16;
      },
      getWidth: () => {
        return 20;
      },
      getVGap: () => {
        return 15;
      },
      getHGap: () => {
        return 80;
      },
      getSide: () => {
        return 'right';
      },
    },
  });

  data.graph.node((node) => {
    return {
      label: node.name,
      labelCfg: {
        position: node.children && node.children.length > 0
            ? 'left'
            : 'right',
        offset: 5,
      },
      style: {
        fill: '#5B8FF9', // 填充颜色
        stroke: '#5B8FF9', // 边框颜色
        lineWidth: 2, // 线宽
        opacity: 0.8, // 不透明度
        size: 40, // 节点大小
        endArrow: {
          path: 'M 0,0 L 10,5 L 10,-5 Z', // 箭头形状
          fill: '#999', // 箭头颜色
        },
      }
    };
  });
  data.graph.data({name: '', children: list.value});
  data.graph.render();
  data.graph.fitView();
}

const props = defineProps({
  outLineList: {
    type: Array,
    default: () => []
  },
  showOutLine: {
    type: Boolean,
    default: false
  }
})

const list = computed(() => {
  return convertToHierarchy(props.outLineList)
})

const convertToHierarchy = (answerList) => {
  const hierarchy = [];
  let currentH1 = null;
  let currentH2 = null;

  answerList.forEach(item => {
    item.name = item.content
    if (item.level === 'h1') {
      currentH1 = {...item, children: []};
      hierarchy.push(currentH1);
      currentH2 = null;
    } else if (item.level === 'h2' && currentH1) {
      currentH2 = {...item, children: []};
      currentH1.children.push(currentH2);
    } else if (item.level === 'h3' && currentH2) {
      currentH2.children = currentH2.children || [];
      currentH2.children.push(item);
    }
  });

  return hierarchy;
}
</script>


<style scoped lang="less">

</style>