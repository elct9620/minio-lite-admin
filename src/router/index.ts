import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import AccessKeyView from '../views/AccessKeyView.vue'
import SiteReplicationView from '../views/SiteReplicationView.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: {
      title: 'Dashboard'
    }
  },
  {
    path: '/access-keys',
    name: 'AccessKeys',
    component: AccessKeyView,
    meta: {
      title: 'Access Keys'
    }
  },
  {
    path: '/site-replication',
    name: 'SiteReplication',
    component: SiteReplicationView,
    meta: {
      title: 'Site Replication'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Update document title based on route
router.beforeEach((to) => {
  if (to.meta?.title) {
    document.title = `${to.meta.title} - MinIO Lite Admin`
  }
})

export default router