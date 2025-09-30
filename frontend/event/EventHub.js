import {EventsOn} from "../wailsjs/runtime";
import {EventNames} from "./events";
import {useDownloadStore} from "@/stores/DownloadStore";

/**
 * 初始化 Wails 事件监听器
 */
export function InitWailsListeners() {
    console.log("Initializing Wails event listeners...");
    const downloadStore = useDownloadStore();
    EventsOn(EventNames.DOWNLOAD_NEW, (task) => {
        console.log('Received new download:', task);
        downloadStore.addNewTask(task);
    });

    EventsOn(EventNames.DOWNLOAD_PROGRESS, (progress) => {
        downloadStore.updateTaskProgress(progress);
    });

    EventsOn(EventNames.DOWNLOAD_COMPLETED, (data) => {
        console.log('Received download completed:', data);
        downloadStore.setTaskCompleted(data);
    });


    EventsOn(EventNames.DOWNLOAD_FAILED, (data) => {
        console.error('Received download failed:', data);
        downloadStore.setTaskFailed(data);
    });
}