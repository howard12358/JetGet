<script setup lang="ts">
import {FlashOutline} from "@vicons/ionicons5";
import {Search} from '@vicons/carbon';
import {type MessageApi, NIcon, useMessage} from "naive-ui";
import {ref} from 'vue';
import {DownloadFile} from "../../wailsjs/go/service/DownloadService";

const message: MessageApi = useMessage();

const downloadUrl = ref<string>('');

const handleDownload = (): void => {
  if (!downloadUrl.value.trim()) {
    message.warning('请输入下载链接');
    return;
  }

  DownloadFile(downloadUrl.value).then((result: string) => {
    message.success(result);
    // 可以在成功后清空输入框
    downloadUrl.value = '';
  }).catch((error: any) => {
    message.error('下载失败: ' + error);
  });
}
</script>

<template>
  <n-flex justify="center" class="search-container">
    <n-input
        v-model:value="downloadUrl"
        placeholder="输入下载链接"
        type="textarea"
        clearable
        :autosize="{minRows: 1, maxRows: 3}"
        class="search-input"
        size="large"
        @keydown.enter="handleDownload"
    >
      <template #prefix>
        <n-icon :component="FlashOutline"/>
      </template>
      <template #suffix>
        <n-button
            strong
            secondary
            type="primary"
            size="small"
            icon-placement="right"
            @click="handleDownload"
        >
          下载
          <template #icon>
            <NIcon>
              <Search/>
            </NIcon>
          </template>
        </n-button>
      </template>
    </n-input>
  </n-flex>
</template>

<style scoped>
.search-container {
  margin-top: 80px;
}

.search-input {
  max-width: 400px;
  width: 100%;
}
</style>