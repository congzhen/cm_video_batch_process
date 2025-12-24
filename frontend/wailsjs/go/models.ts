export namespace process {
	
	export class AppData {
	    outputDirectory: string;
	    cpuThread: number;
	    gpu: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outputDirectory = source["outputDirectory"];
	        this.cpuThread = source["cpuThread"];
	        this.gpu = source["gpu"];
	    }
	}
	export class TranscodeParams {
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
	
	    static createFrom(source: any = {}) {
	        return new TranscodeParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.video_codec = source["video_codec"];
	        this.audio_codec = source["audio_codec"];
	        this.video_height = source["video_height"];
	        this.fps = source["fps"];
	        this.video_bitrate = source["video_bitrate"];
	        this.watermark_content = source["watermark_content"];
	        this.watermark_image = source["watermark_image"];
	        this.watermark_placement = source["watermark_placement"];
	        this.rotate = source["rotate"];
	        this.use_gpu = source["use_gpu"];
	        this.cpu_threads = source["cpu_threads"];
	    }
	}

}

