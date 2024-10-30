import React, { useState } from 'react';
import {
  AppstoreOutlined,
  BarChartOutlined,
  CloudOutlined,
  ShopOutlined,
  TeamOutlined,
  UploadOutlined,
  UserOutlined,
  VideoCameraOutlined,
} from '@ant-design/icons';
import { Layout, Menu, theme } from 'antd';
import Rank from './pages/rank'; // 导入 Rank 组件
import Search from './pages/search'; // 假设你还有这个组件

const { Header, Content, Footer, Sider } = Layout;

const items = [
  { key: '1', icon: <UserOutlined />, label: 'Rank' },
  { key: '2', icon: <VideoCameraOutlined />, label: 'Search' },
  // 可以继续添加其他菜单项
];

const App = () => {
  const [headerComponent, setHeaderComponent] = useState(<Rank />); // 默认显示 Rank 组件

  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  const handleMenuClick = ({ key }) => {
    if (key === '1') {
      setHeaderComponent(<Rank />);
    } else if (key === '2') {
      setHeaderComponent(<Search />);
    }
  };

  return (
      <Layout hasSider>
        <Sider style={{ overflow: 'auto', height: '100vh', position: 'fixed', insetInlineStart: 0, top: 0, bottom: 0, scrollbarWidth: 'thin', scrollbarGutter: 'stable' }}>
          <div className="demo-logo-vertical" />
          <Menu theme="dark" mode="inline" defaultSelectedKeys={['1']} onClick={handleMenuClick} items={items} />
        </Sider>
        <Layout style={{ marginInlineStart: 200 }}>
          <Header style={{ padding: 0, background: colorBgContainer }} />
          <Content style={{ margin: '24px 16px 0', overflow: 'initial' }}>
            <div style={{ padding: 24, textAlign: 'center', background: colorBgContainer, borderRadius: borderRadiusLG }}>
              {headerComponent} {/* 动态渲染的头部组件 */}
            </div>
          </Content>
          <Footer style={{ textAlign: 'center' }}>
            Ant Design ©{new Date().getFullYear()} Created by Ant UED
          </Footer>
        </Layout>
      </Layout>
  );
};

export default App;
