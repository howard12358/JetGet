import {defineStore} from "pinia";
import {ref} from "vue";
import {e, m} from '../../wailsjs/go/models';

export const useDownloadStore = defineStore('download', () => {
    // state: 明确类型为 Record<string, m.DownloadTask>，表示一个以任务ID为键，任务对象为值的集合
    const tasks = ref<Record<string, m.DownloadTaskResp>>({});

    // action: 初始化或批量添加任务
    function initTasks(initialTasks: m.DownloadTaskResp[]) {
        for (const task of initialTasks) {
            tasks.value[task.id] = task;
        }
    }

    // action: 新增下载任务
    function addNewTask(task: m.DownloadTaskResp) {
        tasks.value[task.id] = task;
    }

    // action: 更新下载任务进度
    function updateTaskProgress(progress: e.Progress) {
        const task = tasks.value[progress.id];
        if (task) {
            // 这里可以获得完整的类型提示和安全检查
            task.downloadedSize = progress.downloaded;
            task.totalSize = progress.total;
            task.speed = progress.speed;
            task.status = progress.status;
        }
    }

    // action: 设置下载任务为已完成状态
    function setTaskCompleted(progress: e.Progress) {
        const task = tasks.value[progress.id];
        if (task) {
            task.status = progress.status;
            task.downloadedSize = task.totalSize; // 完成时确保进度为100%
            task.speed = 0;
            task.completedAt = progress.completedAt;
        }
    }

    // action: 设置下载任务为失败状态
    function setTaskFailed(progress: e.Progress) {
        const task = tasks.value[progress.id];
        if (task) {
            task.status = progress.status;
            task.errorMessage = progress.msg;
            task.speed = 0;
        }
    }

    return {
        tasks,
        initTasks,
        addNewTask,
        updateTaskProgress,
        setTaskCompleted,
        setTaskFailed
    }
})