<template>
  <n-form
      ref="formRef"
      :model="setting"
      label-placement="left"
      label-align="left"
      label-width="auto"
      size="small"
      style="margin: 8px"
  >
    <n-form-item label="下载目录">
      <n-input
          v-model:value="setting.downloadDir"
          placeholder="选择下载目录"
          disabled
          style="width: 25%"
      />
      <n-button text type="primary" style="margin-left: 8px" @click="openDir">
        更改
      </n-button>
    </n-form-item>

    <n-form-item label="代理地址">
      <n-input-group>
        <n-input-group-label size="small">http(s)://</n-input-group-label>
        <n-input v-model:value="setting.proxyHost" size="small" :style="{ width: '18%' }" placeholder="主机名或IP"/>
        <n-input-number v-model:value="setting.proxyPort" :style="{ width: '10%' }" placeholder="端口号" :show-button="false" :min="1000">
          <template #prefix>:</template>
        </n-input-number>
      </n-input-group>
    </n-form-item>
    
    <n-form-item>
      <n-button type="primary" :disabled="!isConfigChanged" @click="saveConfig">保存配置</n-button>
    </n-form-item>
  </n-form>
</template>

<script setup>
import {ref, onMounted, watch} from 'vue'
import {ChooseDirectory, GetConfig, SaveConfig} from "../../wailsjs/go/service/SysService";
import {useMessage} from "naive-ui";

const formRef = ref(null)
const setting = ref({
  downloadDir: null,
  proxyHost: null,
  proxyPort: null
})

// 原始配置，用于比较
const originalSetting = ref({
  downloadDir: null,
  proxyHost: null,
  proxyPort: null
})

// 检查配置是否发生变化
const isConfigChanged = ref(false)

// 监听配置变化
watch(setting, () => {
  isConfigChanged.value = checkConfigChanged()
}, { deep: true })

onMounted(() => {
  loadConfig()
})

let message = useMessage();

function loadConfig() {
  GetConfig().then((config) => {
    setting.value.downloadDir = config.DownloadDir
    originalSetting.value.downloadDir = config.DownloadDir
    
    // 解析代理地址
    if (config.Proxy && config.Proxy.startsWith("http")) {
      const url = new URL(config.Proxy)
      setting.value.proxyHost = url.hostname
      setting.value.proxyPort = url.port ? parseInt(url.port) : null
      originalSetting.value.proxyHost = url.hostname
      originalSetting.value.proxyPort = url.port ? parseInt(url.port) : null
    } else {
      originalSetting.value.proxyHost = null
      originalSetting.value.proxyPort = null
    }
    
    // 重置变更状态
    isConfigChanged.value = false
  }).catch((err) => {
    console.error("加载配置失败:", err)
  })
}

// 检查配置是否发生变化
function checkConfigChanged() {
  return (
    setting.value.downloadDir !== originalSetting.value.downloadDir ||
    setting.value.proxyHost !== originalSetting.value.proxyHost ||
    setting.value.proxyPort !== originalSetting.value.proxyPort
  )
}

function openDir() {
  ChooseDirectory().then((dir) => {
    console.log(dir)
    setting.value.downloadDir = dir
  })
}

function saveConfig() {
  // 构造代理URL
  let proxyUrl = ""
  if (setting.value.proxyHost) {
    const port = setting.value.proxyPort || "7890"
    proxyUrl = `http://${setting.value.proxyHost}:${port}`
  }
  
  SaveConfig(setting.value.downloadDir, proxyUrl).then(() => {
    message.success("配置保存成功")
    // 更新原始配置
    originalSetting.value.downloadDir = setting.value.downloadDir
    originalSetting.value.proxyHost = setting.value.proxyHost
    originalSetting.value.proxyPort = setting.value.proxyPort
    // 重置变更状态
    isConfigChanged.value = false
  }).catch((err) => {
    console.error("保存配置失败:", err)
    message.error("保存配置失败")
  })
}
</script>
