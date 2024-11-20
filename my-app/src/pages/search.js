import React, { useState } from 'react';
import {Space, Input, Button, Divider, Table, Spin, Alert, Modal, message, Popconfirm, notification} from 'antd';
import axios from 'axios';
import Cookies from 'js-cookie'; // 用于操作 Cookie
import config from '../conf.js';
import { ExclamationCircleOutlined } from '@ant-design/icons';

// 用户信息表格列配置
const userColumns = [
    {
        title: 'Username',
        dataIndex: 'username',
        key: 'username',
    },
    {
        title: 'Name',
        dataIndex: 'name',
        key: 'name',
    },
    {
        title: 'Location',
        dataIndex: 'location',
        key: 'location',
    },
    {
        title: 'Email',
        dataIndex: 'email',
        key: 'email',
    },
    {
        title: 'Bio',
        dataIndex: 'bio',
        key: 'bio',
    },
    {
        title: 'Followers',
        dataIndex: 'followers',
        key: 'followers',
    },
    {
        title: 'Score',
        dataIndex: 'score',
        key: 'score',
    },
    {
        title: 'Possible Nation',
        dataIndex: 'possible_nation',
        key: 'possible_nation',
    },
    {
        title: 'Confidence Level',
        dataIndex: 'confidence_level',
        key: 'confidence_level',
    },
];

// 仓库信息表格列配置
const repoColumns = [
    {
        title: 'Repository',
        dataIndex: 'repo',
        key: 'repo',
    },
    {
        title: 'Commits',
        dataIndex: 'commits',
        key: 'commits',
    },
    {
        title: 'Forks',
        dataIndex: 'fork',
        key: 'fork',
    },
    {
        title: 'Stars',
        dataIndex: 'star',
        key: 'star',
    },
    {
        title: 'Line Changes',
        dataIndex: 'line_change',
        key: 'line_change',
    },
    {
        title: 'Code Review',
        dataIndex: 'code_review',
        key: 'code_review',
    },
];

const UserReposDisplay = () => {
    const [username, setUsername] = useState('');
    const [userData, setUserData] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [token, setToken] = useState(''); // 用于存储输入的 Token
    const [api, contextHolder] = notification.useNotification();

    const showModal = () => {
        setIsModalOpen(true);
    };

    const handleCancel = () => {
        api.open({
            message: '注意',
            description:
                '默认 token 在流量较高时可能会触发限速，不建议使用。为更稳定的体验，建议配置并使用您自己的 token。' ,
            showProgress: true,
            pauseOnHover:true,
            icon: <ExclamationCircleOutlined style={{ color: '#faad14' }} />,
        });
        // 删除所有的 Cookie
        Cookies.remove('githubToken');
        setIsModalOpen(false);
    };

    // 提交 Token 到后端并存储到 Cookie
    const handleSubmitToken = async () => {
        if (!token) {
            message.error("Token 不能为空");
            return;
        }

        // 简单的 Token 格式检查
        if (token.length < 40) {
            message.error("Token 格式不正确");
            return;
        }

        try {
            const response = await axios.post(`${config.apiBaseUrl}/setToken`, { token });
            if (response.data.message === "Invalid GitHub Token") {
                message.error("无效的 GitHub Token");
                return;
            }
            // Token 验证成功，存储到 Cookie
            Cookies.set('githubToken', token, { expires: 7 }); // 保存到 Cookie，7 天过期
            message.success("Token 设置成功");
            setIsModalOpen(false);
        } catch (err) {
            console.error('Error submitting token:', err);
            // 打印更多错误信息以调试
            message.error('Token 无效、过期');
        }
    };

    const fetchUserData = async () => {
        setLoading(true);  // 开始加载
        setError(null);     // 清除之前的错误

        try {
            const token = Cookies.get('githubToken'); // 从 Cookie 中获取 Token
            if (token) {
                // 验证 Token 是否有效
                const response = await axios.post(`${config.apiBaseUrl}/setToken`, { token });
                if (response.data.message === "Invalid GitHub Token") {
                    throw new Error("无效的 GitHub Token");
                }
            }



            // 获取用户数据
            const userDataResponse = await axios.post(`${config.apiBaseUrl}/userRepos`, {
                username: username, // 传递的 JSON 数据
            },{
                withCredentials: true, // 发送请求时携带 Cookie
            });
            setUserData(userDataResponse.data); // 将返回的数据存储到state中
        } catch (err) {
            console.error('Error fetching user data:', err);
            setError(err.message || '无法获取用户数据，请稍后再试'); // 设置错误信息
        } finally {
            setLoading(false); // 加载结束
        }
    };

    // 用户数据表格数据源
    const userDataSource = userData
        ? [
            {
                key: '1',
                username: userData.username,
                name: userData.name,
                location: userData.location || 'Not provided',
                email: userData.email || 'Not provided',
                bio: userData.bio || 'Not provided',
                followers: userData.followers,
                score: userData.score,
                possible_nation: userData.possible_nation,
                confidence_level: userData.confidence_level,
            },
        ]
        : [];

    // 仓库数据表格数据源
    const repoDataSource = userData && userData.Repos
        ? userData.Repos.map((repo, index) => ({
            key: index.toString(),
            repo: repo.repo,
            commits: repo.commits,
            fork: repo.fork,
            star: repo.star,
            line_change: repo.line_change,
            code_review: repo.code_review,
        }))
        : [];

    return (
        <div>
            {contextHolder}
            <div style={{ textAlign: "left", marginBottom: 20 }}>
                <Button type="primary" onClick={showModal}>
                    Token
                </Button>
                <Modal title="设置 Github Token" open={isModalOpen} onOk={handleSubmitToken} onCancel={handleCancel}
                       cancelText="默认 Token" okText="Submit">
                    <Input
                        placeholder="Enter GitHub Token"
                        value={token}
                        onChange={(e) => setToken(e.target.value)}
                    />
                    <Alert
                        message="隐私声明"
                        description="我们承诺严格保护用户隐私，不保存用户输入的 token，不将其用于商业用途或分享给第三方，
                        仅用于完成必要操作，处理后立即清除，确保数据安全与透明。"
                        type="warning"
                        showIcon
                        style={{ textAlign: "left", marginTop: 10 }}
                    />
                </Modal>
            </div>
            <Space.Compact
                style={{
                    width: '100%',
                }}
            >
                <Input
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="Enter GitHub username"
                />
                <Button type="primary" onClick={fetchUserData} disabled={loading}>
                    Submit
                </Button>
            </Space.Compact>

            {/* 显示加载中的提示 */}
            {loading && (
                <div style={{ marginTop: 20 }}>
                    <Spin tip="加载中..." />
                </div>
            )}

            {/* 显示错误提示 */}
            {error && (
                <div style={{ marginTop: 20 }}>
                    <Alert message={error} type="error" />
                </div>
            )}

            {/* 显示用户信息表格 */}
            {userData && !loading && !error && (
                <div style={{ marginTop: 20 }}>
                    <h3>User Information</h3>
                    <Table columns={userColumns} dataSource={userDataSource} pagination={false} />

                    <Divider />

                    {/* 显示仓库信息表格 */}
                    <h3>Repositories</h3>
                    {repoDataSource.length > 0 ? (
                        <Table columns={repoColumns} dataSource={repoDataSource} pagination={false} />
                    ) : (
                        <p>该用户没有仓库数据。</p>
                    )}
                </div>
            )}
        </div>
    );
};

export default UserReposDisplay;
