<script setup>
import {FlashOutline} from "@vicons/ionicons5";
import {Search} from '@vicons/carbon'
import {NIcon, useMessage} from "naive-ui";
import {ref} from 'vue'
import {DownloadFile} from "../../wailsjs/go/service/DownloadService";

const message = useMessage()
const downloadUrl = ref('')

const handleDownload = () => {
  if (!downloadUrl.value) {
    message.warning('请输入下载链接')
    return
  }

  // 调用后端下载方法
  DownloadFile(downloadUrl.value).then((result) => {
    message.success(result)
  }).catch((error) => {
    message.error('下载失败: ' + error)
  })
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