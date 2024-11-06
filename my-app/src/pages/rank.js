import React, { useEffect, useState } from 'react';
import { Table, Spin, Alert } from 'antd';

const columns = [
    {
        title: 'Rank No',
        dataIndex: 'rankno',
        sorter: false,
        width: '10%',
    },
    {
        title: 'Name',
        dataIndex: 'name',
        sorter: false,
        width: '15%',
    },
    {
        title: 'UserName',
        dataIndex: 'username',
        sorter: false,
        width: '15%',
    },
    {
        title: 'Location',
        dataIndex: 'location',
        filters: [
            { text: 'China', value: 'China' },
        ],
        onFilter: (value, record) => record.location.includes(value),
        width: '15%',
    },
    {
        title: 'Score',
        dataIndex: 'score',
        sorter: false,
        width: '20%',
    },
    {
        title: 'possible_nation',
        dataIndex: 'possible_nation',
        sorter: false,
        width: '15%',
    },
    {
        title: 'confidence_level',
        dataIndex: 'confidence_level',
        sorter: false,
        width: '10%',
    },
];

const Rank = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [pagination, setPagination] = useState({ pageSize: 10, current: 1, total: 0 });
    const [filter, setFilter] = useState({ location: null });

    useEffect(() => {
        const fetchData = async () => {
            const { pageSize, current } = pagination;
            const locationParam = filter.location ? `&location=${filter.location}` : '';

            try {
                const response = await fetch(
                    `http://localhost:8050/user?page_size=${pageSize}&page_number=${current}${locationParam}`
                );
                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error(`Network response was not ok: ${errorText}`);
                }
                const result = await response.json();
                console.log("Fetched data:", result); // 打印获取的数据

                // 更新数据和分页信息
                if (Array.isArray(result.users)) {
                    setData(result.users);
                    setPagination((prevPagination) => ({
                        ...prevPagination,
                        total: result.total, // 从接口返回的总数更新
                    }));
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
    }, [pagination.current, pagination.pageSize, filter]); // 只在 pagination 或 filter 变化时触发

    const handleTableChange = (pagination, filters) => {
        const { current, pageSize } = pagination;

        // 如果筛选条件没有变化则不更新
        setPagination((prevPagination) => ({
            ...prevPagination,
            current,
            pageSize,
        }));

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
            rowKey={(record) => record.id}
            dataSource={data}
            pagination={{
                pageSize: pagination.pageSize,
                current: pagination.current,
                total: pagination.total,
                showSizeChanger: true,
                pageSizeOptions: ['10', '20', '50'],
            }}
            onChange={handleTableChange} // 处理分页和筛选变化
        />
    );
};

export default Rank;
