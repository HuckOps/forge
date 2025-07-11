import { deletePushGateway, getPushGateways } from '@/services/api';
import { PlusOutlined } from '@ant-design/icons';
import ProCard from '@ant-design/pro-card';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { Button, message, Popconfirm } from 'antd';
import { useEffect, useState } from 'react';

export default function () {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [loading, setLoading] = useState(false);

  const [total, setTotal] = useState(0);
  const [pushGateway, setPushGateway] = useState([]);

  const getData = async () => {
    setLoading(true);
    const result = await getPushGateways(
      (currentPage - 1) * pageSize,
      pageSize,
    );
    if (result) {
      console.log(result.data.data);
      setLoading(false);
      setTotal(result.data.total);
      setPushGateway(result.data.data);
    }
  };

  const gotoAdd = () => {
    history.push('/prometheus/pushgateway/add');
  };

  useEffect(() => {
    getData().then();
  }, []);

  return (
    <PageContainer>
      <ProCard>
        <ProTable
          dataSource={pushGateway}
          search={false}
          loading={loading}
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            pageSizeOptions: [10, 20, 50],
            total,

            onChange: (page, pageSize) => {
              setPageSize(pageSize);
              setCurrentPage(page);
            },
          }}
          // @ts-ignore
          toolBarRender={() => (
            <>
              <Button
                key="button"
                icon={<PlusOutlined />}
                type="primary"
                onClick={() => {
                  gotoAdd();
                }}
              >
                新增部署
              </Button>
            </>
          )}
          columns={[
            {
              title: '端口',
              dataIndex: 'port',
            },
            {
              title: '版本',
              dataIndex: 'version',
            },
            {
              title: '机器UUID',
              dataIndex: 'uuid',
            },
            {
              title: '机器IP',
              dataIndex: 'ip',
            },
            {
              title: '操作',
              render: (text, record: any) => (
                <>
                  <Button
                    type="link"
                    onClick={() => {
                      history.push(
                        `/prometheus/pushgateway/add?id=${record?.id}`,
                      );
                    }}
                  >
                    编辑
                  </Button>
                  <Popconfirm
                    title={'确认删除PushGateway吗？'}
                    onConfirm={async () => {
                      const result = await deletePushGateway(record?.id);
                      if (result) {
                        message.success('删除PushGateway成功');
                        await getData();
                      }
                    }}
                  >
                    <Button type={'link'}>删除</Button>
                  </Popconfirm>
                </>
              ),
            },
          ]}
        />
      </ProCard>
    </PageContainer>
  );
}
