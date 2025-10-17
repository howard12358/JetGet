<template>
  <div class="download-page-container">
    <n-tabs
        type="bar"
        animated
        default-value="downloading"
        style="display: flex; flex-direction: column; flex-grow: 1; overflow: hidden;"
        pane-style="flex-grow: 1; overflow: hidden;"
    >
      <n-tab-pane
          name="downloading"
          tab="下载中"
          style="height: 100%; display: flex; flex-direction: column;"
      >
        <n-scrollbar style="flex-grow: 1;">
          <n-list hoverable clickable style="padding-right: 16px; padding-bottom: 20px;">
            <n-list-item v-if="downloadingList.length === 0">
              <n-empty description="暂无下载任务"/>
            </n-list-item>
            <n-list-item v-for="task in downloadingList" :key="task.id">
              <n-thing>
                <template #header>{{ task.fileName }}</template>
                <template #header-extra>
                  <n-tag :type="statusMap[task.status]?.type || 'default'" size="small">
                    {{ statusMap[task.status]?.label || task.status }}
                  </n-tag>
                </template>
                <template #description>
                  <div v-if="task.status === 'failed'" style="color: #d03050">
                    {{ task.errorMessage }}
                  </div>

                  <div v-else style="display: flex; justify-content: space-between; align-items: center;">
                    <span>{{ formatBytes(task.downloadedSize) }} / {{ formatBytes(task.totalSize) }}</span>
                    <span v-show="task.speed !== undefined && !isNaN(task.speed) && task.speed !== -1"
                          style="font-size: 12px; color: #888;">{{ formatBytes(task.speed || 0) }}/s</span>
                  </div>
                </template>

                <n-progress
                    type="line"
                    :percentage="task.totalSize > 0 ? Number(((task.downloadedSize / task.totalSize) * 100).toFixed(2)) : 0"
                    processing
                />
              </n-thing>
            </n-list-item>
          </n-list>
        </n-scrollbar>
      </n-tab-pane>

      <n-tab-pane
          name="downloaded"
          tab="已下载"
          style="height: 100%; display: flex; flex-direction: column;"
      >
        <n-scrollbar style="flex-grow: 1;">
          <n-list hoverable clickable style="padding-right: 16px; padding-bottom: 20px;">
            <n-list-item v-if="completedList.length === 0">
              <n-empty description="暂无下载记录"/>
            </n-list-item>
            <n-list-item v-for="task in completedList" :key="task.id">
              <n-thing :title="task.fileName">
                <template #description>
                  文件大小: {{ formatBytes(task.totalSize) }}
                </template>
                <template #header-extra>
                  <div style="display: flex; align-items: center; gap: 8px;">
                    <n-time :time="new Date(task.completedAt)" format="yyyy-MM-dd HH:mm"/>
                    <n-button text type="error">
                      <template #icon>
                        <Delete/>
                      </template>
                    </n-button>
                  </div>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-scrollbar>
      </n-tab-pane>
    </n-tabs>

  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from "vue";
import {PageDownloadHistory} from "../../wailsjs/go/service/DownloadService";
import {useMessage} from "naive-ui";
import {useDownloadStore} from "@/stores/DownloadStore";
import {TaskStatus} from "@/types/types";
import {Delete} from "@vicons/carbon";

const message = useMessage();
const downloadStore = useDownloadStore();

// --- 状态映射，用于UI展示 ---
const statusMap = {
  downloading: {label: '下载中', type: 'info'},
  pending: {label: '等待中', type: 'warning'},
  completed: {label: '已完成', type: 'success'},
  failed: {label: '失败', type: 'error'},
};

// 动态计算任务列表
const downloadingList = computed(() => {
  return Object.values(downloadStore.tasks)
      .filter(task => task.status !== TaskStatus.StatusCompleted)
      .sort((a, b) => (new Date(b.createdAt)).getTime() - (new Date(a.createdAt)).getTime());
});

const completedList = computed(() => {
  return Object.values(downloadStore.tasks)
      .filter(task => task.status === TaskStatus.StatusCompleted)
      .sort((a, b) => (new Date(b.completedAt)).getTime() - (new Date(a.completedAt)).getTime());
});

// 工具函数：格式化字节
function formatBytes(bytes: number, decimals = 2): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

// 组件挂载时加载历史记录
onMounted(async () => {
  try {
    // 加载未完成的任务
    const downloadingResult = await PageDownloadHistory('downloading', 1, 100);
    if (downloadingResult.list) {
      downloadStore.initTasks(downloadingResult.list);
    }

    // 加载已完成的任务
    const completedResult = await PageDownloadHistory('completed', 1, 100);
    if (completedResult.list) {
      downloadStore.initTasks(completedResult.list);
    }
  } catch (error: any) {
    console.error("加载历史记录失败:", error);
    message.error(`加载历史记录失败: ${error.message || error}`);
  }
});
</script>

<style scoped>
.download-page-container {
  display: flex;
  flex-direction: column;

  /*
     必须使用 calc 来获取一个明确的视口高度，减去父级 layout 的 padding (16px * 2)
  */
  height: calc(100vh - 32px);
}
</style>