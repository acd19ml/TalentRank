import React, { useState } from 'react';
import { Space, Input, Button, Divider, Card, Spin, Alert } from 'antd';
import axios from 'axios';

const UserReposDisplay = () => {
    const [username, setUsername] = useState('');
    const [userData, setUserData] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const fetchUserData = async () => {
        setLoading(true);  // 开始加载
        setError(null);     // 清除之前的错误

        try {
            // 发送 POST 请求，传递 JSON 数据
            const response = await axios.post('http://localhost:8050/userRepos', {
                username: username, // 传递的 JSON 数据
            });
            setUserData(response.data); // 将返回的数据存储到state中
        } catch (err) {
            console.error('Error fetching user data:', err);
            setError('无法获取用户数据，请稍后再试'); // 设置错误信息
        } finally {
            setLoading(false); // 加载结束
        }
    };

    return (
        <div>
            <Space.Compact style={{ width: '100%' }}>
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

            {/* 显示返回的数据 */}
            {userData && !loading && !error && (
                <div style={{ marginTop: 20 }}>
                    <Card title="User Information" style={{ marginBottom: 20 }}>
                        <p><strong>Username:</strong> {userData.username}</p>
                        <p><strong>Name:</strong> {userData.name}</p>
                        <p><strong>Location:</strong> {userData.location}</p>
                        <p><strong>Email:</strong> {userData.email || 'Not provided'}</p>
                        <p><strong>Bio:</strong> {userData.bio || 'Not provided'}</p>
                        <p><strong>Followers:</strong> {userData.followers}</p>
                        <p><strong>Score:</strong> {userData.score}</p>
                        <p><strong>Possible Nation:</strong> {userData.possible_nation}</p>
                        <p><strong>Confidence Level:</strong> {userData.confidence_level}</p>
                    </Card>

                    <Divider />

                    <h3>Repos</h3>
                    {userData.Repos && userData.Repos.length > 0 ? (
                        userData.Repos.map((repo) => (
                            <Card key={repo.id} title={repo.repo} style={{ marginBottom: 20 }}>
                                <p><strong>Commits:</strong> {repo.commits}</p>
                                <p><strong>Forks:</strong> {repo.fork}</p>
                                <p><strong>Stars:</strong> {repo.star}</p>
                                <p><strong>Line Changes:</strong> {repo.line_change}</p>
                                <p><strong>Code Review:</strong> {repo.code_review}</p>
                            </Card>
                        ))
                    ) : (
                        <p>该用户没有仓库数据。</p>
                    )}
                </div>
            )}
        </div>
    );
};

export default UserReposDisplay;
