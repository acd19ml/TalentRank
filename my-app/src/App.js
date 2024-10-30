import React, { useState } from 'react';
import {
  UserOutlined,
  VideoCameraOutlined,
} from '@ant-design/icons';
import { Layout, Menu, theme } from 'antd';
import Rank from './pages/rank'; // 导入 Rank 组件
import Search from './pages/search'; // 假设你还有这个组件

const { Header, Content, Footer, Sider } = Layout;

// 定义菜单项及对应的组件
const menuItems = {
  '1': { icon: <UserOutlined />, label: 'Rank', component: <Rank /> ,text:'贡献分数排名' },
  '2': { icon: <VideoCameraOutlined />, label: 'Search', component: <Search /> ,text:'搜索' },
  // 继续添加其他菜单项
};

const items = Object.keys(menuItems).map(key => ({
  key,
  icon: menuItems[key].icon,
  label: menuItems[key].label,
}));

const App = () => {
  const [headerComponent, setHeaderComponent] = useState(menuItems['1'].component); // 默认显示 Rank 组件
  const [headerLabel, setHeaderLabel] = useState(menuItems['1'].text); // 默认显示的标签

  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  const handleMenuClick = ({ key }) => {
    // 更新组件和标签
    setHeaderComponent(menuItems[key].component);
    setHeaderLabel(menuItems[key].text);
  };

  return (
      <Layout hasSider>
        <Sider style={{ overflow: 'auto', height: '100vh', position: 'fixed', insetInlineStart: 0, top: 0, bottom: 0, scrollbarWidth: 'thin', scrollbarGutter: 'stable' }}>
          <div className="demo-logo-vertical" />
          <Menu theme="dark" mode="inline" defaultSelectedKeys={['1']} onClick={handleMenuClick} items={items} />
        </Sider>
        <Layout style={{ marginInlineStart: 200 }}>
          <Header style={{ padding: 0, background: colorBgContainer }}>
            <h2 style={{ margin: 0, padding: '0 24px' }}>{headerLabel}</h2> {/* 显示当前标签 */}
          </Header>
          <Content style={{ margin: '24px 16px 0', overflow: 'initial' }}>
            <div style={{ padding: 24, textAlign: 'center', background: colorBgContainer, borderRadius: borderRadiusLG }}>
              {headerComponent} {/* 动态渲染的头部组件 */}
            </div>
          </Content>
          <Footer style={{ textAlign: 'center' }}>
            小团体 ©{new Date().getFullYear()} Created by Ant UED
          </Footer>
        </Layout>
      </Layout>
  );
};

export default App;
