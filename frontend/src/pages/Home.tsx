import { useEffect, useState } from 'react';
import { Card, Typography, Space, Tag, Spin, Alert } from 'antd';
import { CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import { checkHealth } from '../services/api';

const { Title, Paragraph, Text } = Typography;

interface HealthData {
  status: string;
  timestamp: string;
  version: string;
}

function Home() {
  const [loading, setLoading] = useState(true);
  const [health, setHealth] = useState<HealthData | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchHealth();
  }, []);

  const fetchHealth = async () => {
    try {
      setLoading(true);
      setError(null);
      const response: any = await checkHealth();
      if (response.code === 200) {
        setHealth(response.data);
      } else {
        setError(response.message || '获取健康状态失败');
      }
    } catch (err) {
      setError('无法连接到后端服务');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Card>
        <Title level={2}>AI Skeleton</Title>
        <Paragraph>
          这是一个面向 AI 辅助开发的全栈项目脚手架，专为中文开发者设计。
        </Paragraph>
        <Paragraph>
          <Text strong>技术栈：</Text>
        </Paragraph>
        <Space wrap>
          <Tag color="blue">React</Tag>
          <Tag color="blue">Vite</Tag>
          <Tag color="blue">Ant Design v5</Tag>
          <Tag color="green">Go</Tag>
          <Tag color="green">Gin</Tag>
          <Tag color="green">GORM</Tag>
        </Space>
      </Card>

      <Card title="后端服务状态">
        {loading ? (
          <Spin tip="检查中..." />
        ) : error ? (
          <Alert
            message="服务异常"
            description={error}
            type="error"
            icon={<CloseCircleOutlined />}
            showIcon
          />
        ) : health ? (
          <Space direction="vertical" size="middle">
            <Space>
              <Text strong>状态：</Text>
              <Tag
                icon={<CheckCircleOutlined />}
                color="success"
              >
                {health.status}
              </Tag>
            </Space>
            <Space>
              <Text strong>版本：</Text>
              <Text>{health.version}</Text>
            </Space>
            <Space>
              <Text strong>时间：</Text>
              <Text type="secondary">{health.timestamp}</Text>
            </Space>
          </Space>
        ) : null}
      </Card>
    </Space>
  );
}

export default Home;
