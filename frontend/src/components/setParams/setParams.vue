<template>
    <div class="set-params-container">
        <el-form :model="videoParams" label-width="auto">
            <div class="block-container">
                <div class="block">
                    <el-form-item label="视频编码">
                        <selectVideoCodec v-model="videoParams.video_codec" :width="props.formWidth">
                        </selectVideoCodec>
                    </el-form-item>
                    <el-form-item label="音频编码">
                        <selectAudioCodec v-model="videoParams.audio_codec" :width="props.formWidth">
                        </selectAudioCodec>
                    </el-form-item>
                    <el-form-item label="视频码率">
                        <selectVideoBitrate v-model="videoParams.video_bitrate" :width="props.formWidth"
                            :allowCreate="true">
                        </selectVideoBitrate>
                    </el-form-item>

                </div>
                <div class="block">
                    <el-form-item label="视频帧率">
                        <selectFps v-model="videoParams.fps" :width="props.formWidth" :allowCreate="true"></selectFps>
                    </el-form-item>
                    <el-form-item label="视频旋转">
                        <selectRotate v-model="videoParams.rotate" :width="props.formWidth"></selectRotate>
                    </el-form-item>
                    <el-form-item label="视频尺寸">
                        <selectVideoHeight v-model="videoParams.video_height" :width="props.formWidth"
                            :allowCreate="true">
                        </selectVideoHeight>
                    </el-form-item>
                </div>
                <div class="block">
                    <el-form-item label="水印文字">
                        <div :style="{ width: props.formWidth }">
                            <el-input v-model="videoParams.watermark_content" placeholder="水印文字"
                                width="100%"></el-input>
                        </div>
                    </el-form-item>
                    <el-form-item label="图片水印">
                        <div :style="{ width: props.formWidth }">
                            <el-input v-model="videoParams.watermark_image">
                                <template #append>
                                    <div class="openWatermarkImageDialog" @click="openWatermarkImageDialogHandle">
                                        <el-icon>
                                            <FolderOpened />
                                        </el-icon>
                                    </div>
                                </template>
                            </el-input>
                        </div>
                    </el-form-item>
                    <el-form-item label="水印位置">
                        <selectWatermarkPlacement v-model="videoParams.watermark_placement" :width="props.formWidth">
                        </selectWatermarkPlacement>
                    </el-form-item>
                </div>
                <div class="block">

                    <el-form-item label="CPU线程">
                        <el-input-number v-model="videoParams.cpu_threads" :min="0" :max="props.cpuThreads"
                            controls-position="right" width="60px" />
                    </el-form-item>
                    <el-form-item>
                        <el-checkbox v-model="videoParams.use_gpu" label="使用GPU加速" :disabled="!props.gpuStatus" />
                    </el-form-item>
                </div>
            </div>
        </el-form>
    </div>
</template>
<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import selectVideoCodec from '../comForm/selectVideoCodec.vue';
import selectAudioCodec from '../comForm/selectAudioCodec.vue';
import selectVideoHeight from '../comForm/selectVideoHeight.vue';
import selectFps from '../comForm/selectFps.vue';
import selectRotate from '../comForm/selectRotate.vue';
import selectWatermarkPlacement from '../comForm/selectWatermarkPlacement.vue';
import selectVideoBitrate from '../comForm/selectVideoBitrate.vue';
import type { videoParams } from '../../datatype/app.datatype';
import { EventsOn_watermarkImageDialog, openWatermarkImageDialog } from '../../process/dialog.process';
const props = defineProps({
    formWidth: {
        type: String,
        default: '220px',
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

const videoParams = ref<videoParams>({
    video_codec: 'copy',
    audio_codec: 'copy',
    video_height: 'copy',
    video_bitrate: 'copy',
    fps: 'copy',
    rotate: 'copy',
    watermark_content: '',
    watermark_image: '',
    watermark_placement: 'top-right',
    use_gpu: false,
    cpu_threads: 0,
});
// 监听 watermarkContent 并过滤非法字符
watch(() => videoParams.value.watermark_content, (newVal) => { // 只允许字母、数字、中文和普通空格
    const filtered = newVal.replace(/[^\w\u4e00-\u9fa5 ]/g, '');
    if (newVal !== filtered) {
        videoParams.value.watermark_content = filtered;
    }
});

const openWatermarkImageDialogHandle = async () => {
    await openWatermarkImageDialog();
}

const getVideoParams = () => {
    return videoParams.value;
};
const setVideoParams = (params: videoParams) => {
    videoParams.value = params;
};
const reset = () => {
    videoParams.value = {
        video_codec: 'copy',
        audio_codec: 'copy',
        video_height: 'copy',
        video_bitrate: 'copy',
        fps: 'copy',
        rotate: 'copy',
        watermark_content: '',
        watermark_image: '',
        watermark_placement: 'top-right',
        use_gpu: false,
        cpu_threads: 0,
    }
}

onMounted(() => {
    EventsOn_watermarkImageDialog((filePath: string) => {
        videoParams.value.watermark_image = filePath;
    })
});

defineExpose({
    getVideoParams,
    setVideoParams,
    reset,
});

</script>
<style lang="scss" scoped>
.set-params-container {
    .block-container {
        display: flex;
        gap: 10px;
        flex-direction: column;

        .block {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            align-items: center;
        }
    }


    :deep(.el-form-item) {
        margin-bottom: 5px;
    }


    :deep(.el-input-group__append) {
        padding: 0;

        .openWatermarkImageDialog {
            padding: 0 10px;
            cursor: pointer;
        }
    }
}
</style>