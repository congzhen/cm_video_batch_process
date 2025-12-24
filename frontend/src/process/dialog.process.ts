import { videoInfo } from "@/datatype/app.datatype";
import { OpenMultipleVideoFilesDialog, OpenDirectoryDialogSetOutput, OpenWatermarkImageDialog } from "../../wailsjs/go/process/App";
import { EventsOn } from "../../wailsjs/runtime";
export const openVideoDialog = async () => {
    return await OpenMultipleVideoFilesDialog();
};

export const EventsOn_filesSelectedMultipleVideoFiles = (callback: (arg0: videoInfo[]) => void) => {
    // 监听文件选择事件
    EventsOn("filesSelectedMultipleVideoFilesSuccess", (videoInfoSlc: videoInfo[]) => {
        callback(videoInfoSlc)
    });
    EventsOn("filesSelectedMultipleVideoFilesError", () => {
        callback([])
    });
    EventsOn("filesSelectedMultipleVideoFilesCancelled", () => {
        callback([])
    });
};

export const openDirectoryDialogSetOutput = async () => {
    return await OpenDirectoryDialogSetOutput();
};

export const EventsOn_directoryDialogSetOutput = (callback: (arg0: string) => void) => {
    // 监听选择事件
    EventsOn("directorySelectedSetOutput", (outputDirectory: string) => {
        callback(outputDirectory)
    });
}

export const openWatermarkImageDialog = async () => {
    return await OpenWatermarkImageDialog();
};
export const EventsOn_watermarkImageDialog = (callback: (arg0: string) => void) => {
    // 监听选择事件
    EventsOn("fileSelectedWatermarkImageSuccess", (watermarkImagePath: string) => {
        callback(watermarkImagePath)
    });
}