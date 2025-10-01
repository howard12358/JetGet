<script setup lang="ts">
import type {MenuOption} from 'naive-ui'
import {NIcon} from "naive-ui";
import {type Component, h, onMounted, ref, type VNodeChild, watch} from "vue";
import {useRoute, useRouter} from "vue-router";
import {InitWailsListeners} from "@/event/EventHub";

const router = useRouter();
const route = useRoute();

function renderIcon(icon: Component): () => VNodeChild {
  return () => h(NIcon, null, {default: () => h(icon)});
}

// 从路由表动态生成菜单，排除 Setting
const mainMenuOptions: MenuOption[] = router.getRoutes()
    .filter(r => r.meta && r.meta.label && r.name !== 'Setting')
    .map(r => ({
      label: r.meta.label as string, // 断言为 string
      key: r.path,
      icon: r.meta.icon ? renderIcon(r.meta.icon as Component) : undefined
    }));

// Setting 菜单单独处理
const settingMenuOptions: MenuOption[] = router.getRoutes()
    .filter(r => r.meta && r.meta.label && r.name === 'Setting')
    .map(r => ({
      label: r.meta.label as string,
      key: r.path,
      icon: r.meta.icon ? renderIcon(r.meta.icon as Component) : undefined
    }));

const selectedKey = ref<string>(route.path);
watch(() => route.path, (newPath: string) => {
  selectedKey.value = newPath;
});

function handleMenuUpdate(key: string): void {
  if (key && key !== route.path) {
    router.push(key);
  }
}

onMounted(() => {
  InitWailsListeners();
});
</script>

<template>
  <div class="app-root">
    <n-message-provider>
      <!-- 整个布局占满视口高度 -->
      <n-layout has-sider style="height: 100vh;">
        <n-layout-sider
            bordered
            collapse-mode="width"
            :collapsed-width="58"
            :native-scrollbar="false"
            collapsed
            style="height:100vh; position: relative;"
        >
          <!-- 主菜单 -->
          <n-menu
              v-model:value="selectedKey"
              :options="mainMenuOptions"
              collapsed
              :collapsed-width="58"
              :collapsed-icon-size="22"
              @update:value="handleMenuUpdate"
              style="padding: 5px 0 80px;"
              class="main-menu"
          />

          <!-- 底部固定的 Setting 菜单 -->
          <div style="position: absolute; bottom: 5px; width: 100%;">
            <n-menu
                v-model:value="selectedKey"
                :options="settingMenuOptions"
                collapsed
                :collapsed-width="58"
                :collapsed-icon-size="22"
                @update:value="handleMenuUpdate"
                class="bottom-menu"
            />
          </div>
        </n-layout-sider>

        <!-- 右侧主区域，同样撑满高度并可滚动 -->
        <n-layout style="height:100vh; padding:16px">
          <n-layout-content>
            <!-- 路由视图：路由内容会在这里渲染 -->
            <router-view/>
          </n-layout-content>
        </n-layout>
      </n-layout>
    </n-message-provider>
  </div>

</template>
