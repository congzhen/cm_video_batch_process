import { OnFileDrop } from "../../wailsjs/runtime";

export const EventsOn_OnFileDrop = async () => {
    await OnFileDrop(async (x: number, y: number, paths: string[]) => {

    }, true);
};