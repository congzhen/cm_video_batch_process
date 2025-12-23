<template>
    <el-select v-model="selectVal" :clearable="props.clearable" :style="{ width: props.width }" placeholder="视频码率"
        @change="changeHandle" @clear="handleClear" :multiple="props.multiple">
        <el-option v-for="item, index in dataset.videoBitrate" :key="index" :label="getFormatFileSize(item)"
            :value="item"></el-option>
    </el-select>
</template>
<script setup lang="ts">
import dataset from '@/assets/dataset';
import { formatFileSize } from '@/assets/dataConversion';
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

const getFormatFileSize = (bytes: string): string => {
    if (bytes === 'copy') {
        return 'copy';
    }
    return formatFileSize(parseInt(bytes));
}

const changeHandle = () => {
    emit('change', selectVal.value || '')
}
const handleClear = () => {
    if (props.multiple) {
        selectVal.value = [];
    } else {
        selectVal.value = '';
    }
}
</script>
