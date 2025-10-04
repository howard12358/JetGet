import {createRouter, createWebHistory, RouteRecordRaw} from 'vue-router';
import type {Component} from 'vue';
import {ArrowDownload16Regular} from '@vicons/fluent';
import {Search, Settings} from '@vicons/carbon';

// [最佳实践] 通过模块扩展来为 RouteMeta 添加自定义属性的类型定义
// 这样做可以让你在整个项目中（例如导航守卫或组件内）访问 route.meta 时获得类型提示
declare module 'vue-router' {
    interface RouteMeta {
        label: string;
        icon: Component;
    }
}

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/search'
    },
    {
        path: '/search',
        name: 'Search',
        component: () => import('@/pages/Search.vue'),
        meta: {
            label: '搜索',
            icon: Search
        }
    },
    {
        path: '/download',
        name: 'Download',
        component: () => import('@/pages/Download.vue'),
        meta: {
            label: '下载',
            icon: ArrowDownload16Regular
        }
    },
    {
        path: '/setting',
        name: 'Setting',
        component: () => import('@/pages/Setting.vue'),
        meta: {
            label: '设置',
            icon: Settings
        }
    }
];

const router = createRouter({
    // @ts-expect-error
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
});

export default router;