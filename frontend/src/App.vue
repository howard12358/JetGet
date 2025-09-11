<script setup>
import {NIcon} from "naive-ui";
import {h, ref, watch} from "vue";
import {useRoute, useRouter} from "vue-router";

const router = useRouter();
const route = useRoute();

function renderIcon(icon) {
  return () => h(NIcon, null, {default: () => h(icon)});
}

// 从路由表动态生成菜单
const menuOptions = router.getRoutes()
    // 只取那些带 meta.label 的
    .filter(r => r.meta && r.meta.label)
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
</script>

<template>
  <div class="app-root">
    <n-message-provider>
      <!-- 整个布局占满视口高度 -->
      <n-layout has-sider style="height: 100vh;">
        <n-layout-sider
            bordered
            collapse-mode="width"
            :collapsed-width="64"
            :native-scrollbar="false"
            collapsed
            style="height:100vh; display:flex; flex-direction:column;"
        >
          <!-- 用 v-model:value 双向绑定选中项，并监听 update:value 以跳转 -->
          <n-menu
              v-model:value="selectedKey"
              :options="menuOptions"
              collapsed
              :collapsed-width="64"
              :collapsed-icon-size="22"
              @update:value="handleMenuUpdate"
              style="flex: 1 1 auto;"
          />
        </n-layout-sider>

        <!-- 右侧主区域，同样撑满高度并可滚动 -->
        <n-layout style="height:100vh;">
          <n-layout-content style="height:100vh; overflow:auto; padding:16px;">
            <!-- 路由视图：路由内容会在这里渲染 -->
            <router-view/>
          </n-layout-content>
        </n-layout>
      </n-layout>
    </n-message-provider>
  </div>

</template>
