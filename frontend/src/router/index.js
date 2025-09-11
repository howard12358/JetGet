import {createRouter, createWebHistory} from 'vue-router';
import { ArrowDownload16Regular } from '@vicons/fluent'
import { Settings } from '@vicons/carbon'

const routes = [
    {
        path: '/',
        redirect: '/download'
    },
    {
        path: '/download',
        name: 'Download',
        component: () => import('@/pages/Download.vue'),
        meta: {
            label: 'Download',
            icon: ArrowDownload16Regular
        }
    },
    {
        path: '/setting',
        name: 'Setting',
        component: () => import('@/pages/Setting.vue'),
        meta: {
            label: 'Setting',
            icon: Settings
        }
    }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
});

export default router;