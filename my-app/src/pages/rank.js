import React, { useEffect, useState } from 'react';
import { Table, Spin, Alert } from 'antd';

const columns = [
    {
        title: 'Rank No',
        dataIndex: 'rankno',
        sorter: true,
        width: '20%',
    },
    {
        title: 'Name',
        dataIndex: 'name',
        sorter: true,
        width: '30%',
    },
    {
        title: 'Location',
        dataIndex: 'location',
        filters: [
            { text: 'China', value: 'China' },
            { text: 'USA', value: 'USA' },
            // { text: 'Location C', value: 'Location C' },
            // 可以根据你的数据添加更多位置
        ],
        onFilter: (value, record) => record.location.includes(value), // 筛选逻辑
        width: '30%',
    },
    {
        title: 'Score',
        dataIndex: 'score',
        sorter: true,
        width: '20%',
    },
];

const Rank = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/rankings');
                if (!response.ok) {
                    const errorText = await response.text(); // 获取详细的错误信息
                    throw new Error(`Network response was not ok: ${errorText}`);
                }
                const result = await response.json();
                setData(result);
            } catch (err) {
                console.error("Fetch error:", err); // 日志输出
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    if (loading) {
        return <Spin tip="Loading..." />;
    }

    if (error) {
        return <Alert message="Error" description={error} type="error" />;
    }

    return (
        <Table
            columns={columns}
            rowKey={(record) => record.rankno} // Assuming rankno is unique
            dataSource={data}
        />
    );
};

export default Rank;
