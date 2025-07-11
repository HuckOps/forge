import Operate from '@/pages/Nodes/operate';
import { getNodes } from '@/services/api';
import {
  PageContainer,
  ProFormInstance,
  ProTable,
} from '@ant-design/pro-components';
import { Tag } from 'antd';
import React, { useEffect, useRef, useState } from 'react';

export default function () {
  const ref = useRef<ProFormInstance>();

  const [total, setTotal] = useState<number>(0);
  const [nodes, setNodes] = useState([]);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(20);
  const [loading, setLoading] = useState<boolean>(false);

  const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([]);

  const getData = async () => {
    setLoading(true);
    console.log(ref.current?.getFieldsValue());
    const result = await getNodes(
      (currentPage - 1) * pageSize,
      pageSize,
      ref.current?.getFieldsValue(),
    );
    if (result) {
      setNodes(result.data.data);
      setTotal(result.data.total);
      setLoading(false);
      return;
    }
  };

  useEffect(() => {
    getData().then();
  }, [pageSize, currentPage]);

  // @ts-ignore
  // @ts-ignore
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
        rowSelection={{
          align: 'left',
          selectedRowKeys: selectedRowKeys,
          onChange: (newSelectedRowKeys: React.Key[]) =>
            //   @ts-ignore
            setSelectedRowKeys(newSelectedRowKeys),
        }}
        rowKey={'id'}
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
        toolbar={{
          search: (
            <Operate
              SelectedRowKeys={selectedRowKeys}
              resetSelectRowKeys={() => {
                setSelectedRowKeys([]);
              }}
            />
          ),
        }}
      ></ProTable>
    </PageContainer>
  );
}
