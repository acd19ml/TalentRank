import React from 'react';
import './App.css'; // 引入样式文件

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      title: 'TalentRank',
      width: 200,
      collapsed: false,
      currentPage: <div>欢迎使用小团体的TalentRank!</div>,
      menus: [
        {
          text: "Forms",
          iconCls: "fa fa-wpforms",
          children: [
            { text: "Form Element", page: <div>Form Element Content</div> },
            { text: "Wizard", page: <div>Wizard Content</div> },
            { text: "File Upload", page: <div>File Upload Content</div> }
          ]
        },
        {
          text: "Mail",
          iconCls: "fa fa-at",
          children: [
            { text: "Inbox", page: <div>Inbox Content</div> },
            { text: "Sent", page: <div>Sent Content</div> },
            { text: "Trash", page: <div>Trash Content</div> }
          ]
        },
        {
          text: "Layout",
          iconCls: "fa fa-table",
          children: [
            { text: "Panel", page: <div>Panel Content</div> },
            { text: "Accordion", page: <div>Accordion Content</div> },
            { text: "Tabs", page: <div>Tabs Content</div> }
          ]
        }
      ]
    };
  }

  toggle() {
    const { collapsed } = this.state;
    this.setState({
      collapsed: !collapsed,
      width: collapsed ? 200 : 50
    });
  }

  handleItemClick(item) {
    this.setState({ currentPage: item.page });
  }

  render() {
    const { menus, title, width, collapsed, currentPage } = this.state;

    return (
        <div style={{ display: 'flex', height: '100vh' }}>
          <div className="sidebar-body" style={{ width: width, backgroundColor: 'rgb(51, 51, 51)', color: '#fff' }}>
            <div className="sidebar-header" style={{ padding: '10px', textAlign: 'center' }}>
              <h3>{title}</h3>
            </div>
          </div>
          <div style={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
            <div className="main-header" style={{ backgroundColor: '#f8f9fa', padding: '10px', display: 'flex', alignItems: 'center' }}>
              <span className="main-toggle fa fa-bars" onClick={this.toggle.bind(this)} style={{ cursor: 'pointer' }}></span>
              <h2 style={{ marginLeft: '10px' }}>{title}</h2>
            </div>
            <div className="main-body" style={{ padding: '20px', flexGrow: 1, backgroundColor: '#f0f0f0' }}>
              {currentPage}
            </div>
          </div>
        </div>
    );
  }
}

export default App;
