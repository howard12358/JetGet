import {EventsOn} from "../../wailsjs/runtime";
import {useDownloadStore} from "@/stores/DownloadStore";
import {e, m} from "../../wailsjs/go/models";
import {EventNames} from "@/types/types";

/**
 * 初始化 Wails 事件监听器
 */
export function InitWailsListeners() {
    console.log("Initializing Wails event listeners...");
    const downloadStore = useDownloadStore();

    // 监听新任务事件
    EventsOn(EventNames.DOWNLOAD_NEW, (data: any) => {
        console.log('Received new download:', data);
        const newTask = m.DownloadTaskResp.createFrom(data);
        downloadStore.addNewTask(newTask);
    });

    // 监听进度事件
    EventsOn(EventNames.DOWNLOAD_PROGRESS, (data: any) => {
        const progress = e.Progress.createFrom(data);
        downloadStore.updateTaskProgress(progress);
    });

    // 监听任务完成事件
    EventsOn(EventNames.DOWNLOAD_COMPLETED, (data: any) => {
        console.log('Received download completed:', data);
        const completedData = e.Progress.createFrom(data);
        downloadStore.setTaskCompleted(completedData);
    });

    // 监听任务失败事件
    EventsOn(EventNames.DOWNLOAD_FAILED, (data: any) => {
        console.error('Received download failed:', data);
        const failedData = e.Progress.createFrom(data);
        downloadStore.setTaskFailed(failedData);
    });
}