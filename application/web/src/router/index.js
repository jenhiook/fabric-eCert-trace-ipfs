import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * 常量路由（无需权限）
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },
  {
    path: '/home',
    component: Layout,
    children: [{
      path: '',
      name: 'Home',
      component: () => import('@/views/home/index'),
      meta: { title: '系统首页', icon: 'el-icon-s-home' }
    }]
  },
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/uplink',
    component: Layout,
    children: [{
      path: '',
      name: 'Uplink',
      component: () => import('@/views/uplink/index'),
      meta: { title: '溯源信息录入', icon: 'el-icon-edit-outline', roles: ['个人用户'] }
    }]
  },
  {
    path: '/govt-audit',
    component: Layout,
    children: [{
      path: '',
      name: 'GovtAudit',
      component: () => import('@/views/govt-audit/index'),
      meta: { title: '证照审核', icon: 'el-icon-check', roles: ['政务部门'] }
    }]
  },
  {
    path: '/enterprise',
    component: Layout,
    children: [{
      path: '',
      name: 'Enterprise',
      component: () => import('@/views/enterprise/index'),
      meta: { title: '证照使用备案', icon: 'el-icon-s-order', roles: ['企业组织'] }
    }]
  },
  {
    path: '/tech-support',
    component: Layout,
    children: [{
      path: '',
      name: 'TechSupport',
      component: () => import('@/views/tech-support/index'),
      meta: { title: '证照核验', icon: 'el-icon-circle-check', roles: ['技术支撑实体'] }
    }]
  },
  {
    path: '/trace',
    component: Layout,
    children: [{
      path: '',
      name: 'Trace',
      component: () => import('@/views/trace/index'),
      meta: { title: '溯源查询', icon: 'el-icon-search' }
    }]
  },
  {
    path: '/dashboard',
    component: Layout,
    children: [{
      path: '',
      name: 'Dashboard',
      component: () => import('@/views/dashboard/index'),
      meta: { title: '可视化分析', icon: 'el-icon-data-line',roles: ['个人用户'] }
    }]
  },
  {
    path: '/my-certs',
    component: Layout,
    children: [{
      path: '',
      name: 'MyCerts',
      component: () => import('@/views/my-certs/index'),
      meta: { title: '我的证照', icon: 'el-icon-s-claim', roles: ['个人用户'] }
    }]
  },
  {
    path: '/timeline',
    component: Layout,
    children: [{
      path: '',
      name: 'TimelineOnly',
      component: () => import('@/views/trace/TimelineOnly.vue'),
      meta: { title: '证照时间轴', icon: 'el-icon-time', roles: ['个人用户'] }
    }]
  },
  // ==================== 政府端新增路由 ====================
  {
    path: '/gov-trace',
    component: Layout,
    children: [{
      path: '',
      name: 'GovTrace',
      component: () => import('@/views/gov/Trace.vue'),
      meta: { title: '政府溯源', icon: 'el-icon-s-management', roles: ['政务部门'] }
    }]
  },
  // 404 must be at the end
  { path: '*', redirect: '/404', hidden: true },
  {
    path: '/gov-dashboard',
    component: Layout,
    children: [{
      path: '',
      name: 'GovDashboard',
      component: () => import('@/views/gov/Dashboard.vue'),
      meta: { title: '政府可视化', icon: 'el-icon-data-analysis', roles: ['政务部门'] }
    }]
  },
  {
    path: '/enterprise-dashboard',
    component: Layout,
    children: [{
      path: '',
      name: 'EnterpriseDashboard',
      component: () => import('@/views/enterprise/Dashboard.vue'),
      meta: { title: '企业可视化', icon: 'el-icon-data-analysis', roles: ['企业组织'] }
   }]
  },
  {
    path: '/enterprise-trace',
    component: Layout,
    children: [{
      path: '',
      name: 'EnterpriseTrace',
      component: () => import('@/views/enterprise/Trace.vue'),
      meta: { title: '企业溯源', icon: 'el-icon-search', roles: ['企业组织'] }
   }]
  },
  {
    path: '/enterprise-supplier',
    component: Layout,
    children: [{
      path: '',
      name: 'EnterpriseSupplier',
      component: () => import('@/views/enterprise/SupplierLink.vue'),
      meta: { title: '供应商关联', icon: 'el-icon-connection', roles: ['企业组织'] }
   }]
  },
  {
    path: '/enterprise-event',
    component: Layout,
    children: [{
      path: '',
      name: 'EnterpriseEvent',
      component: () => import('@/views/enterprise/ComplianceEvent.vue'),
      meta: { title: '合规事件', icon: 'el-icon-warning', roles: ['企业组织'] }
   }]
  },
  {
    path: '/tech-verify',
    component: Layout,
    children: [{
      path: '',
      name: 'TechVerify',
      component: () => import('@/views/tech/Verify.vue'),
      meta: { title: '核验工具', icon: 'el-icon-check', roles: ['技术支撑实体'] }
   }]
  },
  {
    path: '/tech-dashboard',
    component: Layout,
    children: [{
      path: '',
      name: 'TechDashboard',
      component: () => import('@/views/tech/Dashboard.vue'),
      meta: { title: '技术可视化', icon: 'el-icon-data-analysis', roles: ['技术支撑实体'] }
   }]
  }
]

const createRouter = () => new Router({
  // mode: 'history', // 根据需要开启
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// 重置路由（用于动态权限）
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher
}

export default router
