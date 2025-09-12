<script setup>
import {FlashOutline} from "@vicons/ionicons5";
import { Search } from '@vicons/carbon'
import {NIcon} from "naive-ui";
</script>

<template>
  <n-flex justify="center" class="search-container">
    <n-input placeholder="输入链接"
             type="textarea"
             clearable
             :autosize="{minRows: 1, maxRows: 3}"
             class="search-input"
             size="large"
    >
      <template #prefix>
        <n-icon :component="FlashOutline"/>
      </template>
      <template #suffix>
        <n-button strong secondary type="primary" size="small" icon-placement="right">
          搜索
          <template #icon>
            <NIcon>
              <Search />
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