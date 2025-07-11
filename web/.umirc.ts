import { defineConfig } from '@umijs/max';

export default defineConfig({
  antd: {},
  access: {},
  model: {},
  initialState: {},
  request: {},
  layout: {
    title: '@umijs/max',
  },
  routes: [
    {
      path: '/',
      redirect: '/home',
    },
    {
      name: '首页',
      path: '/home',
      component: './Home',
    },
    {
      name: '节点',
      path: '/nodes',
      routes: [
        {
          name: '节点管理',
          path: '/nodes/list',
          component: 'Nodes',
        },
        {
          name: '节点标签',
          path: '/nodes/labels',
          component: 'Labels',
        },
        {
          path: '/nodes/labels/add',
          component: 'Labels/add.tsx',
        },
        {
          name: '标签信息',
          path: '/nodes/labels/:id',
          component: 'Labels/detail.tsx',
          hideInMenu: true,
        },
      ],
    },
    {
      name: 'Prometheus管理',
      path: '/prometheus',
      routes: [
        {
          name: 'PushGateway管理',
          path: '/prometheus/pushgateway',
          component: 'PushGateway',
        },
        {
          name: 'PushGateway管理',
          path: '/prometheus/pushgateway/add',
          component: 'PushGateway/add',
          hideInMenu: true,
        },
        {
          name: '联邦节点管理',
          path: '/prometheus/federation',
          component: 'Federation',
        },
        {
          name: '联邦节点管理',
          path: '/prometheus/federation/add',
          component: 'Federation/add',
          hideInMenu: true,
        },
      ],
    },
  ],
  npmClient: 'yarn',
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
});
