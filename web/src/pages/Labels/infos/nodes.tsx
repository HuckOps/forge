import { getLabelNodes } from '@/services/api';
import { ProTable } from '@ant-design/pro-components';
import { Tag } from 'antd';
import { useEffect, useState } from 'react';

interface Props {
  labelID?: string;
}
export default function (props: Props) {
  const [total, setTotal] = useState<number>(0);
  const [nodes, setNodes] = useState<API.Node[]>([]);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(20);
  const [loading, setLoading] = useState<boolean>(false);
  const getData = async () => {
    setLoading(true);
    const result = await getLabelNodes(
      props?.labelID,
      (currentPage - 1) * pageSize,
      pageSize,
    );
    if (result) {
      console.log(result.data.data);
      setNodes(result.data.data);
      setTotal(result.data.total);
      setLoading(false);
      return;
    }
  };
  useEffect(() => {
    getData().then();
  }, [pageSize, currentPage]);
  return (
    <>
      <ProTable
        dataSource={nodes}
        loading={loading}
        onSubmit={() => getData()}
        search={false}
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
      ></ProTable>
    </>
  );
}
