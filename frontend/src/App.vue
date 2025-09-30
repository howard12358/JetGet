<script setup>
import {NIcon} from "naive-ui";
import {h, onMounted, ref, watch} from "vue";
import {useRoute, useRouter} from "vue-router";
import {InitWailsListeners} from "../event/EventHub";

const router = useRouter();
const route = useRoute();

function renderIcon(icon) {
  return () => h(NIcon, null, {default: () => h(icon)});
}

// 从路由表动态生成菜单，排除 Setting
const mainMenuOptions = router.getRoutes()
    .filter(r => r.meta && r.meta.label && r.name !== 'Setting')
    .map(r => ({
      label: r.meta.label,
      key: r.path,
      icon: r.meta.icon ? renderIcon(r.meta.icon) : undefined
    }))

// Setting 菜单单独处理
const settingMenuOptions = router.getRoutes()
    .filter(r => r.meta && r.meta.label && r.name === 'Setting')
    .map(r => ({
      label: r.meta.label,
      key: r.path,
      icon: r.meta.icon ? renderIcon(r.meta.icon) : undefined
    }))

// 选中项和路由同步
const selectedKey = ref(route.path)
watch(() => route.path, (p) => {
  selectedKey.value = p
})

function handleMenuUpdate(val) {
  const key = Array.isArray(val) ? val[0] : val
  if (key && key !== route.path) {
    router.push(key)
  }
}

onMounted(()=>{
  InitWailsListeners()
})
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
