import React, { useState } from 'react';
import { Input, Button, Space, message } from 'antd';
import * as echarts from 'echarts';

const Echart2 = () => {
    const [inputValue, setInputValue] = useState('');

    const handleSubmit = async () => {
        console.log(inputValue); // 在控制台打印输入的值

        try {
            const response = await fetch('http://localhost:8080/api/getUserData', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username: inputValue }),
            });

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            const data = await response.json();
            loadChartData(data); // 数据加载后渲染图表
        } catch (error) {
            message.error('请求失败: ' + error.message);
        }
    };

    const loadChartData = (data) => {
        const chartDom = document.getElementById('main');
        const myChart = echarts.init(chartDom);
        myChart.showLoading();

        // 检查数据格式是否正确
        if (!data || !data.name || !Array.isArray(data.children) || data.children.length === 0) {
            message.error('未查询到此用户');
            myChart.hideLoading();
            return;
        }

        // 设置图表选项
        myChart.setOption({
            tooltip: {
                trigger: 'item',
                triggerOn: 'mousemove'
            },
            series: [
                {
                    type: 'tree',
                    data: [data],
                    top: '1%',
                    left: '7%',
                    bottom: '1%',
                    right: '20%',
                    symbolSize: 7,
                    label: {
                        position: 'left',
                        verticalAlign: 'middle',
                        align: 'right',
                        fontSize: 14, // 增大字体
                        formatter: (params) => {
                            const maxLength = 10; // 设置最大字符长度
                            return params.name.length > maxLength ? params.name.slice(0, maxLength) + '...' : params.name;
                        }
                    },
                    leaves: {
                        label: {
                            position: 'right',
                            verticalAlign: 'middle',
                            align: 'left',
                            fontSize: 14, // 增大字体
                            formatter: (params) => {
                                const maxLength = 10; // 设置最大字符长度
                                return params.name.length > maxLength ? params.name.slice(0, maxLength) + '...' : params.name;
                            }
                        }
                    },
                    emphasis: {
                        focus: 'descendant'
                    },
                    expandAndCollapse: true,
                    animationDuration: 550,
                    animationDurationUpdate: 750
                }
            ]
        });

        myChart.hideLoading();
    };

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
            <Space.Compact style={{ width: '100%', marginBottom: '20px' }}>
                <Input
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)} // 更新状态
                />
                <Button type="primary" onClick={handleSubmit}>
                    Submit
                </Button>
            </Space.Compact>
            <div id="main" style={{ flex: 1, minHeight: '400px', marginTop: '20px' }}></div> {/* 图表容器自适应，最小高度为 400px */}
        </div>
    );
};

export default Echart2;
