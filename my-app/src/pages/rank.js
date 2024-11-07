import React, { useEffect, useState } from 'react';
import { Table, Spin, Alert, Popconfirm, message } from 'antd';

const Rank = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [pagination, setPagination] = useState({ pageSize: 10, current: 1, total: 0 });
    const [filter, setFilter] = useState({ location: null });

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

    // 删除用户的函数
    const handleDelete = async (id) => {
        try {
            const response = await fetch(`http://localhost:8050/userRepos${id}`, {
                method: 'DELETE',
            });

            if (response.ok) {
                message.success('User deleted successfully');
                setLoading(true); // 显示 loading 状态
                setPagination((prev) => ({
                    ...prev,
                    current: 1, // 删除后刷新第一页
                }));
                // 手动触发数据重新获取
                fetchData();
            } else {
                const errorText = await response.text();
                throw new Error(`Failed to delete: ${errorText}`);
            }
        } catch (error) {
            message.error(`Error: ${error.message}`);
        }
    };


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
            width: '15%',
        },
        {
            title: 'score',
            dataIndex: 'score',
            sorter: false,
            width: '20%',
        },
        {
            title: 'possible_nation',
            dataIndex: 'possible_nation',
            filters: [
                { text: 'China', value: 'China' },
            ],
            onFilter: (value, record) => record.location.includes(value),
            sorter: false,
            width: '15%',
        },
        {
            title: 'confidence_level',
            dataIndex: 'confidence_level',
            sorter: false,
            width: '10%',
        },
        {
            title: 'operation',
            dataIndex: 'operation',
            render: (_, record) => (
                <Popconfirm
                    title="Are you sure to delete this user?"
                    onConfirm={() => handleDelete(record.id)} // 调用删除函数
                >
                    <a>Delete</a>
                </Popconfirm>
            ),
        },
    ];

    useEffect(() => {

        fetchData();
    }, [pagination.current, pagination.pageSize, filter]);

    const handleTableChange = (pagination, filters) => {
        const { current, pageSize } = pagination;
        setPagination((prevPagination) => ({
            ...prevPagination,
            current,
            pageSize,
        }));

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
