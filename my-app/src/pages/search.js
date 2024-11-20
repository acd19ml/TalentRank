import React, { useState } from 'react';
import { Space, Input, Button, Divider, Table, Spin, Alert, Modal, message } from 'antd';
import axios from 'axios';
import Cookies from 'js-cookie'; // 用于操作 Cookie
import config from '../conf.js';

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

    const showModal = () => {
        setIsModalOpen(true);
    };

    const handleOk = () => {
        setIsModalOpen(false);
    };

    const handleCancel = () => {
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
            if (!token) {
                throw new Error("未设置 GitHub Token");
            }

            // 验证 Token 是否有效
            const response = await axios.post(`${config.apiBaseUrl}/setToken`, { token });
            if (response.data.message === "Invalid GitHub Token") {
                throw new Error("无效的 GitHub Token");
            }

            // 获取用户数据
            const userDataResponse = await axios.post(`${config.apiBaseUrl}/userRepos`, {
                username: username, // 传递的 JSON 数据
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
            <div style={{ textAlign: "left", marginBottom: 20 }}>
                <Button type="primary" onClick={showModal}>
                    Token
                </Button>
                <Modal title="Set GitHub Token" open={isModalOpen} onOk={handleSubmitToken} onCancel={handleCancel} >
                    <Input
                        placeholder="Enter GitHub Token"
                        value={token}
                        onChange={(e) => setToken(e.target.value)}
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
