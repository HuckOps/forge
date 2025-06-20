import { getNodes } from '@/services/api';
import {
  PageContainer,
  ProFormInstance,
  ProTable,
} from '@ant-design/pro-components';
import { Tag } from 'antd';
import { useEffect, useRef, useState } from 'react';

export default function () {
  const ref = useRef<ProFormInstance>();

  const [total, setTotal] = useState<number>(0);
  const [nodes, setNodes] = useState([]);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(20);
  const [loading, setLoading] = useState<boolean>(false);

  const getData = async () => {
    setLoading(true);
    console.log(ref.current?.getFieldsValue());
    const result = await getNodes(
      (currentPage - 1) * pageSize,
      pageSize,
      ref.current?.getFieldsValue(),
    );
    if (result) {
      setNodes(result.data.nodes);
      setTotal(result.data.total);
      setLoading(false);
      return;
    }
  };

  useEffect(() => {
    getData().then();
  }, [pageSize, currentPage]);

  return (
    <PageContainer
      header={{
        title: '节点管理',
      }}
    >
      <ProTable
        formRef={ref}
        dataSource={nodes}
        loading={loading}
        onSubmit={() => getData()}
        columns={[
          { title: 'UUID', dataIndex: 'uuid' },
          {
            title: '首次注册时间',
            dataIndex: 'created_at',
            hideInSearch: true,
          },
          {
            title: '主机名',
            dataIndex: 'hostname',
          },
          {
            title: 'IP',
            dataIndex: 'ip',
          },
          {
            title: '心跳上报时间',
            dataIndex: 'heartbeat',
            hideInSearch: true,
          },
          {
            title: '心跳状态',
            dataIndex: 'heartBeatStatus',
            valueType: 'select',
            valueEnum: {
              true: '正常',
              false: '异常',
            },
            render: (text) =>
              text ? (
                <Tag color={'green'}>正常</Tag>
              ) : (
                <Tag color={'red'}>异常</Tag>
              ),
          },
        ]}
        pagination={{
          total,
          current: currentPage,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50'],
          pageSize: pageSize,
          onChange: (page, size) => {
            setPageSize(size);
            setCurrentPage(page);
          },
        }}
      ></ProTable>
    </PageContainer>
  );
}
