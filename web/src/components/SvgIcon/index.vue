<template>
    <div v-if="isOnlineSvg" :style="{ '--svg-icon-url': `url(${iconClass})` }" class="svg-icon svg-icon-online"
        :class="svgClass" />
    <svg v-else :class="svgClass">
        <use :xlink:href="iconName" />
    </svg>
</template>
<script setup>

import { defineProps, computed } from 'vue'

const props = defineProps({
    iconClass: {
        type: String,
        required: true
    },
    className: {
        type: String,
        default: ''
    }
})

const isOnlineSvg = computed(() => {
    return /^(https?:)/.test(props.iconClass) || /^(http?:)/.test(props.iconClass)
})

const iconName = computed(() => {
    return `#icon-${props.iconClass}`
})

const svgClass = computed(() => {
    if (props.className) {
        return 'svg-icon ' + props.className
    } else {
        return 'svg-icon'
    }
})

</script>

<style scoped>
.svg-icon {
    /* width: 28px;
  height: 28px; */
    vertical-align: -0.15em;
    fill: currentColor;
    overflow: hidden;
}

.svg-external-icon {
    background-color: currentColor;
    mask-size: cover !important;
    display: inline-block;
}

.svg-icon-online {
    background-color: currentColor;
    fill: currentColor;
    overflow: hidden;
    mask-image: var(--svg-icon-url);
    -webkit-mask-image: var(--svg-icon-url);
    mask-size: cover !important;
    -webkit-mask-size: cover;
    display: inline-block;
}
</style>
