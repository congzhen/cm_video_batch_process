export interface AppData {
    outputDirectory: string;
    cpuThread: number;
    gpu: boolean;
}

export interface videoInfo {
    id: string;
    name: string,
    path: string,
    thumbnail: string,
    size: number,
    duration: number,
    width: number,
    height: number,
    fps: number,
    bitrate: number,
    video_codec: string,
    audio_codec: string,
    video_bitrate: number,
    audio_bitrate: number,

}

export interface videoInfoHasParams extends videoInfo {
    outputSetParams: null | videoParams,
    progress: number,
}

export interface videoParams {
    video_codec: string;
    audio_codec: string;
    video_height: string;
    fps: string;
    video_bitrate: string;
    watermark_content: string;
    watermark_image: string;
    watermark_placement: string;
    rotate: string;
    use_gpu: boolean;
    cpu_threads: number;
}