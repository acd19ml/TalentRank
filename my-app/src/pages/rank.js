import React, { useEffect, useState } from 'react';
import { Table, Spin, Alert } from 'antd';

const columns = [
    {
        title: 'Rank No',
        dataIndex: 'rankno', // 这里可以根据你的需求调整
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
    const [pagination, setPagination] = useState({ pageSize: 10, pageNumber: 1 });
    const [filter, setFilter] = useState({ location: null });

    useEffect(() => {
        const fetchData = async () => {
            const { pageSize, pageNumber } = pagination;
            const locationParam = filter.location ? `&location=${filter.location}` : '';

            try {
                const response = await fetch(`http://localhost:8050/user?page_size=${pageSize}&page_number=${pageNumber}${locationParam}`);
                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error(`Network response was not ok: ${errorText}`);
                }
                const result = await response.json();
                console.log("Fetched data:", result); // 打印获取的数据

                // 提取 users 数组并设置数据
                if (Array.isArray(result.users)) {
                    setData(result.users);
                } else {
                    throw new Error("Expected an array from the API");
                }
            } catch (err) {
                console.error("Fetch error:", err);
                setError(err.message || "An error occurred");
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [pagination, filter]); // 依赖于 pagination 和 filter

    const handleTableChange = (pagination, filters) => {
        const { current: pageNumber, pageSize } = pagination;
        setPagination({ pageSize, pageNumber });

        // 更新筛选条件
        setFilter({ location: filters.location ? filters.location[0] : null });
    };

    if (loading) {
        return <Spin tip="Loading..." />;
    }

    if (error) {
        return <Alert message="Error" description={error} type="error" />;
    }

    return (
        <Table
            columns={columns}
            rowKey={(record) => record.id} // 假设 id 是唯一的
            dataSource={data}
            pagination={{ pageSize: pagination.pageSize, current: pagination.pageNumber }} // 设置分页
            onChange={handleTableChange} // 处理分页和筛选变化
        />
    );
};

export default Rank;
