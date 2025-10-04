<template>
  <n-form
      ref="formRef"
      :model="setting"
      label-placement="left"
      label-align="left"
      label-width="auto"
      size="small"
      style="margin: 10px"
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
        <n-input-number v-model:value="setting.proxyPort" :style="{ width: '10%' }" placeholder="端口号"
                        :show-button="false" :min="1000">
          <template #prefix>:</template>
        </n-input-number>
      </n-input-group>
    </n-form-item>

    <n-form-item>
      <n-button type="primary" :disabled="!isConfigChanged" @click="saveConfig">保存配置</n-button>
    </n-form-item>
  </n-form>
</template>

<script setup lang="ts">
import {onMounted, ref, watch} from 'vue'
import {ChooseDirectory, GetConfig, SaveConfig} from "../../wailsjs/go/service/SysService";
import {type FormInst, type MessageApi, useMessage} from "naive-ui";
import {m} from "../../wailsjs/go/models";

// 定义表单数据的接口 (Interface)
interface SettingForm {
  downloadDir: string | null;
  proxyHost: string | null;
  proxyPort: number | null;
}

// 为 ref 添加类型
const formRef = ref<FormInst | null>(null);
const setting = ref<SettingForm>({
  downloadDir: null,
  proxyHost: null,
  proxyPort: null
});
// 原始配置，用于比较
const originalSetting = ref<SettingForm>({
  downloadDir: null,
  proxyHost: null,
  proxyPort: null
});

// 检查配置是否发生变化
const isConfigChanged = ref<boolean>(false);

// 监听配置变化
watch(setting, () => {
  isConfigChanged.value = checkConfigChanged();
}, {deep: true});

onMounted(() => {
  loadConfig();
});

// 为 useMessage 添加类型
const message: MessageApi = useMessage();

function loadConfig() {
  // 为 GetConfig 的返回结果添加类型
  GetConfig().then((config: m.SysConfig) => {
    setting.value.downloadDir = config.DownloadDir;
    originalSetting.value.downloadDir = config.DownloadDir;

    // 解析代理地址
    if (config.Proxy && config.Proxy.startsWith("http")) {
      try {
        const url = new URL(config.Proxy);
        setting.value.proxyHost = url.hostname;
        setting.value.proxyPort = url.port ? parseInt(url.port, 10) : null;
        originalSetting.value.proxyHost = url.hostname;
        originalSetting.value.proxyPort = url.port ? parseInt(url.port, 10) : null;
      } catch (e) {
        console.error("解析代理URL失败:", e);
        // 如果解析失败，则重置代理设置
        setting.value.proxyHost = null;
        setting.value.proxyPort = null;
        originalSetting.value.proxyHost = null;
        originalSetting.value.proxyPort = null;
      }
    } else {
      // 如果没有代理或格式不正确，则重置
      setting.value.proxyHost = null;
      setting.value.proxyPort = null;
      originalSetting.value.proxyHost = null;
      originalSetting.value.proxyPort = null;
    }

    // 重置变更状态
    isConfigChanged.value = false;
  }).catch((err: any) => {
    message.error("加载配置失败");
    console.error("加载配置失败:", err);
  });
}

// 检查配置是否发生变化
function checkConfigChanged(): boolean {
  return (
      setting.value.downloadDir !== originalSetting.value.downloadDir ||
      setting.value.proxyHost !== originalSetting.value.proxyHost ||
      setting.value.proxyPort !== originalSetting.value.proxyPort
  );
}

function openDir() {
  ChooseDirectory().then((dir: string) => {
    if (dir) { // 用户可能取消选择，这时 dir 为空
      setting.value.downloadDir = dir;
    }
  }).catch((err: any) => {
    console.error("选择目录失败:", err);
  });
}

function saveConfig() {
  // 构造代理URL
  let proxyUrl = "";
  if (setting.value.proxyHost && setting.value.proxyHost.trim() !== '') {
    // 如果没有端口号，可以提供一个默认值，或者根据业务需求决定是否允许为空
    const port = setting.value.proxyPort || 7890; // 假设默认端口是7890
    proxyUrl = `http://${setting.value.proxyHost}:${port}`;
  }

  // 确保 downloadDir 不为 null
  const downloadDir = setting.value.downloadDir || "";

  SaveConfig(downloadDir, proxyUrl).then(() => {
    message.success("配置保存成功");
    // 更新原始配置
    originalSetting.value.downloadDir = setting.value.downloadDir;
    originalSetting.value.proxyHost = setting.value.proxyHost;
    originalSetting.value.proxyPort = setting.value.proxyPort;
    // 重置变更状态
    isConfigChanged.value = false;
  }).catch((err: any) => {
    console.error("保存配置失败:", err);
    message.error("保存配置失败");
  });
}
</script>
