import { Layout } from 'antd';
import Home from './pages/Home';

const { Header, Content, Footer } = Layout;

function App() {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#fff', padding: '0 50px' }}>
        <h1 style={{ margin: 0, lineHeight: '64px' }}>AI Skeleton</h1>
      </Header>
      <Content style={{ padding: '50px' }}>
        <div style={{ maxWidth: 1200, margin: '0 auto' }}>
          <Home />
        </div>
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        AI Skeleton ©2026 Created for AI-Assisted Development
      </Footer>
    </Layout>
  );
}

export default App;
