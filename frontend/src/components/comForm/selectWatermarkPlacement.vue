<template>
    <el-select v-model="selectVal" :clearable="props.clearable" :style="{ width: props.width }" placeholder="水印位置"
        @change="changeHandle" @clear="handleClear" :multiple="props.multiple">
        <el-option v-for="item, index in dataset.watermarkPlacement" :key="index" :label="getLabelCn(item)"
            :value="item"></el-option>
    </el-select>
</template>
<script setup lang="ts">
import dataset from '@/assets/dataset';
const selectVal = defineModel<string | string[]>({ type: [String, Array], default: "" as string | string[] });
const props = defineProps({
    width: {
        type: String,
        default: '100%',
    },
    multiple: {
        type: Boolean,
        default: false
    },
    clearable: {
        type: Boolean,
        default: false
    }
})
const emit = defineEmits(['change'])

const changeHandle = () => {
    emit('change', selectVal.value || '')
}

const getLabelCn = (test: string) => {
    switch (test) {
        case 'top-right':
            return test + ' (右上角)';
            break;
        case 'random':
            return test + ' (四角随机)';
            break;
        case 'horizontal':
            return test + ' (水平横移)';
            break;
        case 'diagonal':
            return test + ' (对角线)';
            break;
        case 'bounce':
            return test + ' (弹跳)';
            break;
        case 'spiral':
            return test + ' (螺旋运动)';

        default:
            return 'Error';
    }
}

const handleClear = () => {
    if (props.multiple) {
        selectVal.value = [];
    } else {
        selectVal.value = '';
    }
}
</script>
