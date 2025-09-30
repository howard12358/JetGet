import {defineStore} from "pinia";
import {ref} from "vue";

export const useDownloadStore = defineStore('download', () => {
    // state: 使用一个对象来存储任务，key 为任务 ID，方便快速查找和更新
    const tasks = ref({});

    // 添加一个新任务
    function addNewTask(task) {
        tasks.value[task.id] = { ...task };
    }

    // 更新任务进度
    function updateTaskProgress(progress) {
        const task = tasks.value[progress.id];
        if (task) {
            task.downloadedSize = progress.downloaded;
            task.totalSize = progress.total;
            task.speed = progress.speed;
            task.status = progress.status;
        }
    }

    // 设置任务完成
    function setTaskCompleted(payload) {
        const task = tasks.value[payload.id];
        if (task) {
            task.status = payload.status;
            task.downloadedSize = task.totalSize;
            task.speed = 0;
        }
    }

    // 设置任务失败
    function setTaskFailed(payload) {
        const task = tasks.value[payload.id];
        if (task) {
            task.status = payload.status;
            task.errorMessage = payload.msg;
            task.speed = 0;
        }
    }

    return {
        tasks,
        addNewTask,
        updateTaskProgress,
        setTaskCompleted,
        setTaskFailed
    }
})