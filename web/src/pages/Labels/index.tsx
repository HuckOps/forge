import { getLabels } from '@/services/api';
import { PlusOutlined } from '@ant-design/icons';
import {
  PageContainer,
  ProFormInstance,
  ProTable,
} from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { Button, Popconfirm } from 'antd';
import { useEffect, useRef, useState } from 'react';

export default function () {
  const ref = useRef<ProFormInstance>();
  const [total, setTotal] = useState<number>(0);
  const [labels, setLabels] = useState([]);

  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const getData = async () => {
    const result = await getLabels(
      (currentPage - 1) * pageSize,
      pageSize,
      ref.current?.getFieldsValue(),
    );
    if (result) {
      setLabels(result.data.data);
      setTotal(result.data.total);
    }
  };

  useEffect(() => {
    getData().then();
  }, [pageSize, currentPage]);

  // @ts-ignore
  return (
    <PageContainer
      header={{
        title: '节点标签',
      }}
    >
      <ProTable
        dataSource={labels}
        onSubmit={getData}
        columns={[
          { title: '名称', dataIndex: 'name' },
          { title: '代码', dataIndex: 'code' },
          { title: '描述', dataIndex: 'description' },
          { title: '创建时间', dataIndex: 'created_at' },
          { title: '更新时间', dataIndex: 'updated_at' },
          {
            title: '操作',
            render: (text, record) => (
              <>
                <Button
                  type={'link'}
                  onClick={() => {
                    // @ts-ignore
                    history.push(`/nodes/labels/${record?.id}`);
                  }}
                >
                  详细
                </Button>
                <Popconfirm title={'确认删除'}>
                  <Button type={'link'}>删除</Button>
                </Popconfirm>
              </>
            ),
          },
        ]}
        pagination={{
          current: currentPage,
          pageSize,
          pageSizeOptions: [10, 20, 50],
          onChange: (page, pageSize) => {
            setPageSize(pageSize);
            setCurrentPage(page);
          },
          total,
        }}
        // @ts-ignore
        toolBarRender={() => (
          <Button
            key="button"
            icon={<PlusOutlined />}
            type="primary"
            onClick={() => history.push('/nodes/labels/add')}
          >
            新增
          </Button>
        )}
      />
    </PageContainer>
  );
}
