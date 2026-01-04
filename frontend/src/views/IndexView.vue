<template>
    <div class="index-container">
        <div class="toolbar">
            <el-button type="primary" icon="Plus" plain @click="openVideoDialogHandle">选择视频</el-button>
            <el-button type="danger" icon="Delete" plain @click="clearHandle">清空列表</el-button>
            <el-button type="info" icon="Refresh" plain @click="resetListHandle">重置列表</el-button>
        </div>
        <div class="video-list">
            <el-table :data="videoList" height="100%" v-loading="loading" empty-text="未选择视频" style="width: 100%">
                <el-table-column label="预览" width="160">
                    <template #default="scope">
                        <el-image class="thumbnail" :src="scope.row.thumbnail" fit="cover">
                        </el-image>
                    </template>
                </el-table-column>
                <el-table-column label="源视频">
                    <template #default="scope">
                        <div class="input-video-info">
                            <div class="video-title">{{ scope.row.name }}</div>
                            <div class="video-tag">
                                <el-tag type="info">{{ scope.row.width + '×' + scope.row.height }}</el-tag>
                                <el-tag type="info">{{ formatFileSize(scope.row.size) }}</el-tag>
                                <el-tag type="info">{{ formatDuration(scope.row.duration) }}</el-tag>
                                <el-tag type="info">{{ scope.row.fps }} fps</el-tag>
                                <el-tag type="info">{{ formatFileSize(scope.row.video_bitrate) }}</el-tag>
                                <el-tag type="info">{{ scope.row.video_codec }}</el-tag>
                                <el-tag type="info">{{ scope.row.audio_codec }}</el-tag>
                            </div>
                            <div class="video-path" :title="scope.row.path">
                                {{ scope.row.path }}
                            </div>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="输出视频">
                    <template #default="scope">
                        <div class="output-video-info">
                            <div class="video-progress">
                                <el-progress :text-inside="true" :stroke-width="17" :percentage="scope.row.progress" />
                            </div>
                            <div class="video-tag" v-if="scope.row.outputSetParams">
                                <el-tag type="info" v-for="item in getSutputSetParams(scope.row.outputSetParams)">
                                    {{ item }}
                                </el-tag>
                            </div>
                            <div class="video-tag" v-else>
                                <el-tag type="primary" effect="light">通用设置</el-tag>
                            </div>
                            <div class="transcodeVideoSuccess" v-if="scope.row.transcodeVideoInfo">
                                <div class="video-tag">
                                    <el-tag type="success">{{ scope.row.transcodeVideoInfo.width + '×' +
                                        scope.row.transcodeVideoInfo.height }}</el-tag>
                                    <el-tag type="success">{{ formatFileSize(scope.row.transcodeVideoInfo.size)
                                        }}</el-tag>
                                    <el-tag type="success">{{ formatDuration(scope.row.transcodeVideoInfo.duration)
                                        }}</el-tag>
                                    <el-tag type="success">{{ scope.row.transcodeVideoInfo.fps }} fps</el-tag>
                                    <el-tag type="success">{{ formatFileSize(scope.row.transcodeVideoInfo.video_bitrate)
                                        }}</el-tag>
                                    <el-tag type="success">{{ scope.row.transcodeVideoInfo.video_codec }}</el-tag>
                                    <el-tag type="success">{{ scope.row.transcodeVideoInfo.audio_codec }}</el-tag>
                                </div>
                            </div>
                        </div>
                    </template>

                </el-table-column>
                <el-table-column label="-" width="110">
                    <template #default="scope">
                        <div class="opt-btn">
                            <div class="opt-btn-item">
                                <el-button type="primary" icon="Setting" plain size="small" title="设置参数"
                                    @click="setParamsDialogHandle(scope.row)" />
                            </div>
                            <div class="opt-btn-item">
                                <el-button type="danger" icon="Delete" plain size="small" title="删除"
                                    @click="deleteVideoHandle(scope.$index)" />
                            </div>
                            <div class="opt-btn-item" v-if="scope.row.transcodeVideoInfo">
                                <el-button type="success" icon="VideoPlay" plain size="small" title="播放"
                                    @click="playTranscodeVideoHandel(scope.row.transcodeVideoInfo.path)" />
                            </div>
                            <div class="opt-btn-item" v-if="scope.row.transcodeVideoInfo">
                                <el-button type="info" icon="Refresh" plain size="small" title="重置"
                                    @click="resetHandel(scope.row)" />
                            </div>
                        </div>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div class="set-params">
            <setParams ref="setParamsRef" :gpu-status="appData?.gpu" :cpu-threads="appData?.cpuThread"></setParams>
        </div>
        <div class="bottom-toolbar">
            <div class="outpath">
                <el-text type="info">输出地址: {{ appData?.outputDirectory }}</el-text>
                <el-link type="primary" @click="openDirectoryDialogSetOutput">选择</el-link>
                <el-link type="primary" @click="openOutputDirectory">打开</el-link>
            </div>
            <div class="btns">
                <div class="show-number">{{ progressCompletedQuantity_C }}/{{ videoList.length }}</div>
                <el-button type="primary" plain @click="startHandle">开始执行</el-button>
            </div>
        </div>
    </div>
    <setParamsDialog ref="setParamsDialogRef" :gpu-status="appData?.gpu" :cpu-threads="appData?.cpuThread">
    </setParamsDialog>
</template>
<script setup lang="ts">
import type { AppData, videoInfo, videoInfoHasParams, videoParams } from '@/datatype/app.datatype';
import { formatFileSize, formatDuration } from '@/assets/dataConversion'
import setParams from '@/components/setParams/setParams.vue';
import { onMounted, ref, computed } from 'vue';
import { EventsOn_filesSelectedMultipleVideoFiles, openVideoDialog, openDirectoryDialogSetOutput, EventsOn_directoryDialogSetOutput } from '@/process/dialog.process'
import { EventsOn_Loading, EventsOn_videoTranscodeProcessor, EventsOn_videoTranscodeSuccess, getAppData, openOutputDirectory, openTranscodeVideo, transcode } from '@/process/app.process'
import setParamsDialog from '@/components/setParams/setParamsDialog.vue';
import { ElMessage } from 'element-plus';
import { EventsOn_OnFileDrop } from '@/process/dragAndDrop.process'

const loading = ref(false)
const setParamsDialogRef = ref<InstanceType<typeof setParamsDialog>>();
const setParamsRef = ref<InstanceType<typeof setParams>>();
const videoList = ref<videoInfoHasParams[]>([])
const appData = ref<AppData>()


const progressCompletedQuantity_C = computed(() => {
    return videoList.value.filter(item => item.progress == 100).length
})

const getSutputSetParams = (params: videoParams) => {
    const arr = []
    if (params.video_codec != 'copy') {
        arr.push('视频编码: ' + params.video_codec)
    }
    if (params.audio_codec != 'copy') {
        arr.push('音频编码: ' + params.audio_codec)
    }
    if (params.fps != 'copy') {
        arr.push('帧率: ' + params.fps)
    }
    if (params.video_bitrate != 'copy') {
        arr.push('视频码率: ' + formatFileSize(parseInt(params.video_bitrate)))
    }
    if (params.video_height != 'copy') {
        arr.push('视频高度: ' + params.video_height)
    }
    if (params.rotate != 'copy') {
        arr.push('旋转: ' + params.rotate + '°')
    }
    if (params.watermark_content != '') {
        arr.push('水印文字: ' + params.watermark_content)
        arr.push('水印位置: ' + params.watermark_placement)
    }
    if (params.use_gpu) {
        arr.push('使用GPU')
    }
    if (params.cpu_threads != 0) {
        arr.push('CPU线程数: ' + params.cpu_threads)
    }
    return arr
}


const openVideoDialogHandle = async () => {
    loading.value = true;
    await openVideoDialog()
}
const clearHandle = () => {
    videoList.value = []
}
const resetListHandle = () => {
    for (const videoInfoHasParams of videoList.value) {
        videoInfoHasParams.progress = 0
        videoInfoHasParams.outputSetParams = null
        videoInfoHasParams.transcodeVideoInfo = null
    }
}
const setParamsDialogHandle = (videoInfoHasParams: videoInfoHasParams) => {
    setParamsDialogRef.value?.open(videoInfoHasParams.name, videoInfoHasParams.outputSetParams, (params: null | videoParams) => {
        console.log(params)
        videoInfoHasParams.outputSetParams = params
    })
}

const deleteVideoHandle = (index: number) => {
    videoList.value.splice(index, 1)
}
const playTranscodeVideoHandel = async (path: string) => {
    await openTranscodeVideo(path)
}

const resetHandel = (videoInfoHasParams: videoInfoHasParams) => {
    videoInfoHasParams.progress = 0
    videoInfoHasParams.outputSetParams = null
    videoInfoHasParams.transcodeVideoInfo = null
}

const startHandle = async () => {
    if (setParamsRef.value) {
        for (const videoInfoHasParams of videoList.value) {
            if (videoInfoHasParams.progress == 100) {
                continue
            }
            const params = videoInfoHasParams.outputSetParams || setParamsRef.value.getVideoParams();
            const result = await transcode(videoInfoHasParams.id, videoInfoHasParams.path, params)
            if (result == 'OK') {
                ElMessage({
                    showClose: true,
                    message: videoInfoHasParams.name + ' 转码完成',
                    type: 'success',
                    duration: 10000,
                });
            } else {
                ElMessage({
                    showClose: true,
                    message: videoInfoHasParams.name + ' 转码失败: ' + result,
                    type: 'error',
                    duration: 10000,
                });
            }
        }
    }
};

onMounted(async () => {
    appData.value = await getAppData()
    EventsOn_Loading((isLoading: boolean) => {
        loading.value = isLoading;
    });
    EventsOn_filesSelectedMultipleVideoFiles((videoInfoSlc: videoInfo[]) => {
        console.log(videoInfoSlc);
        videoList.value.push(...videoInfoSlc.filter(video =>
            !videoList.value.some(existingVideo => existingVideo.path === video.path)
        ).map(videoInfo => {
            return {
                ...videoInfo,
                outputSetParams: null,
                transcodeVideoInfo: null,
                progress: 0
            }
        }))
        loading.value = false;
    })
    EventsOn_directoryDialogSetOutput((directory: string) => {
        if (appData.value) {
            appData.value.outputDirectory = directory
        }
    })
    EventsOn_videoTranscodeProcessor((id: string, progress: number, currentTime: string | 'completed') => {
        for (let i = 0; i < videoList.value.length; i++) {
            if (videoList.value[i].id == id) {
                if (currentTime == 'completed') {
                    videoList.value[i].progress = 100
                } else {
                    videoList.value[i].progress = parseFloat(progress.toFixed(2))
                }
                break
            }
        }
    })
    EventsOn_videoTranscodeSuccess((transcodeVideoInfo: videoInfo) => {
        for (let i = 0; i < videoList.value.length; i++) {
            if (videoList.value[i].id == transcodeVideoInfo.id) {
                videoList.value[i].transcodeVideoInfo = transcodeVideoInfo
                break
            }
        }
    })
    EventsOn_OnFileDrop();
})
</script>
<style lang="scss" scoped>
.index-container {
    width: calc(100% - 20px);
    height: calc(100% - 20px);
    padding: 10px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    overflow: hidden;

    .toolbar {
        flex-shrink: 0;
    }

    .video-list {
        flex: 1;
        overflow: hidden;
        border-top: 1px solid #ddd;
        border-left: 1px solid #ddd;
        border-right: 1px solid #ddd;

        :deep(.el-table__body-wrapper) {
            .cell {
                min-height: 80px;
                overflow: hidden;
                display: flex;
            }
        }

        .thumbnail {
            width: 100%;
            height: 80px;
        }


        .video-tag {
            display: flex;
            flex-wrap: wrap;
            gap: 2px;
        }

        .input-video-info {
            width: 100%;
            height: 100%;
            display: flex;
            flex-direction: column;
            gap: 5px;

            .video-title {
                font-weight: 700;
            }


            .video-path {
                font-size: 11px;
                color: #666;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
        }

        .output-video-info {
            width: 100%;
            height: 100%;
            display: flex;
            flex-direction: column;
            gap: 5px;
        }

        .opt-btn {
            display: flex;
            flex-wrap: wrap; // 允许换行
            width: 100%;
            height: 100%;
            gap: 5px;
            align-content: flex-start; // 顶部对齐


        }
    }

    .set-params {
        flex-shrink: 0;
        display: flex;
    }

    .bottom-toolbar {
        flex-shrink: 0;
        display: flex;
        justify-content: space-between;
        gap: 10px;
        align-items: center;

        .outpath {
            display: flex;
            align-items: center;
            gap: 5px;

            .el-link {
                flex-shrink: 0;
            }
        }

        .btns {
            flex-shrink: 0;
            display: flex;
            gap: 10px;
            align-items: center;

            .show-number {
                flex-shrink: 0;
                display: flex;
                font-size: 12px;
                align-items: center;
            }

            .el-button {
                flex-shrink: 0;
            }
        }

    }
}
</style>