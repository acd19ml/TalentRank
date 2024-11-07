import React, { useState } from 'react';
import {
  UserOutlined,
  VideoCameraOutlined,
  MailOutlined,
} from '@ant-design/icons';
import { Layout, Menu, theme } from 'antd';
import Rank from './pages/rank';
import Search from './pages/search';
import Echart1 from './pages/echart1';
import Echart2 from './pages/echart2';
import Information from './pages/Information';

const { Header, Content, Footer, Sider } = Layout;

const menuItems = {
  '1': { icon: <UserOutlined />, label: 'Rank', component: <Rank />, text: '贡献分数排名' },
  '2': { icon: <VideoCameraOutlined />, label: 'Search', component: <Search />, text: '搜索未记录的用户' },
  '4': { icon: <VideoCameraOutlined />, label: 'information', component: <Information />, text: '信息' },
  '3': {
    icon: <MailOutlined />,
    label: 'echarts',
    children: [
      { key: '31', label: '地区图', component: <Echart1 />, text: 'user地区饼状图' }, // 更新为 Echart1 组件
      { key: '32', label: 'user数据', component: <Echart2 />, text: 'user数据树状图' },
    ],
  },
};

const items = Object.keys(menuItems).map(key => ({
  key,
  icon: menuItems[key].icon,
  label: menuItems[key].label,
  children: menuItems[key].children ? menuItems[key].children.map(child => ({
    key: child.key,
    label: child.label,
  })) : null,
}));

const App = () => {
  const [headerComponent, setHeaderComponent] = useState(menuItems['1'].component);
  const [headerLabel, setHeaderLabel] = useState(menuItems['1'].text);

  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  const handleMenuClick = ({ key }) => {
    const item = Object.values(menuItems).flatMap(menuItem =>
        menuItem.children ? menuItem.children.find(child => child.key === key) : null
    ).find(Boolean);

    if (item) {
      setHeaderComponent(item.component);
      setHeaderLabel(item.text);
    } else {
      setHeaderComponent(menuItems[key]?.component || null);
      setHeaderLabel(menuItems[key]?.text || '');
    }
  };

  return (
      <Layout hasSider>
        <Sider style={{ overflow: 'auto', height: '100vh', position: 'fixed', insetInlineStart: 0, top: 0, bottom: 0, scrollbarWidth: 'thin', scrollbarGutter: 'stable' }}>
          <div className="demo-logo-vertical" />
          <Menu theme="dark" mode="inline" defaultSelectedKeys={['1']} onClick={handleMenuClick} items={items} />
        </Sider>
        <Layout style={{ marginInlineStart: 200 }}>
          <Header style={{ padding: 0, background: colorBgContainer }}>
            <h2 style={{ margin: 0, padding: '0 24px' }}>{headerLabel}</h2>
          </Header>
          <Content style={{ margin: '24px 16px 0', overflow: 'initial' }}>
            <div style={{ padding: 24, textAlign: 'center', background: colorBgContainer, borderRadius: borderRadiusLG }}>
              {headerComponent}
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
