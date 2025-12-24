<template>
    <dialogCommon ref="dialogCommonRef" width="860" :title="title_C" btnSubmitTitle="确定" @submit="submitHandle">
        <setParams class="set-params-dialog" ref="setParamsRef" :form-width="props.formWidth"
            :gpu-status="props.gpuStatus" :cpu-threads="props.cpuThreads"></setParams>
    </dialogCommon>
</template>
<script setup lang="ts">
import type { videoParams } from '@/datatype/app.datatype';
import { computed, ref } from 'vue';
import dialogCommon from '../comDialog/dialog-common.vue';
import setParams from './setParams.vue';
const props = defineProps({
    formWidth: {
        type: String,
        default: '200px',
    },
    gpuStatus: {
        type: Boolean,
        default: false,
    },
    cpuThreads: {
        type: Number,
        default: 0,
    },
});
const dialogCommonRef = ref();
const setParamsRef = ref<InstanceType<typeof setParams>>();
const name = ref('');

let callback: ((params: null | videoParams) => void)
const title_C = computed(() => {
    return `${name.value}`;
});

const submitHandle = () => {
    let params: null | videoParams = null;
    if (setParamsRef.value != null) {
        params = setParamsRef.value.getVideoParams();
    }
    callback(params);
    dialogCommonRef.value.close();
};

const open = (_name: string, _videoInfoHasParams: videoParams | null, _callback: (params: null | videoParams) => void) => {
    name.value = _name;
    callback = _callback;
    if (_videoInfoHasParams != null) {
        setParamsRef.value?.setVideoParams(_videoInfoHasParams)
    } else {
        setParamsRef.value?.reset()
    }
    dialogCommonRef.value.open();
};

defineExpose({
    open,
});


</script>
<style scoped lang="scss"></style>