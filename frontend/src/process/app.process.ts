import { videoParams } from "@/datatype/app.datatype";
import { AppData, OpenOutputDirectory, Transcode } from "../../wailsjs/go/process/App";
import { EventsOn } from "../../wailsjs/runtime";

export const getAppData = async () => {
    return await AppData();
};

export const openOutputDirectory = async () => {
    await OpenOutputDirectory();
};

export const transcode = async (id: string, path: string, params: videoParams): Promise<string> => {
    return await Transcode(id, path, params);
};

export const EventsOn_videoTranscodeProcessor = (callback: (arg0: string, arg1: number, arg2: string | 'completed') => void) => {
    // 监听视频转码进度
    EventsOn("videoTranscodeProcessor", (id: string, process: number, currentTime: string | 'completed') => {
        callback(id, process, currentTime)
    });
};
