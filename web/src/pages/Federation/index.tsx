import { getFederations } from '@/services/api';
import ProCard from '@ant-design/pro-card';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Button } from 'antd';
import { useEffect, useState } from 'react';

export default function () {
  const [total, setTotal] = useState(0);
  const [federation, setFederation] = useState([]);
  const getData = async () => {
    const result = await getFederations();
    // console.log(result);
    if (result) {
      setTotal(result.data.total);
      setFederation(result.data.data);
    }
  };

  useEffect(() => {
    getData();
  }, []);

  return (
    <PageContainer>
      <ProCard>
        <ProTable
          dataSource={federation}
          pagination={{
            total,
          }}
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
              render: (text, record) => (
                <>
                  <Button type={'link'}></Button>
                </>
              ),
            },
          ]}
        ></ProTable>
      </ProCard>
    </PageContainer>
  );
}
