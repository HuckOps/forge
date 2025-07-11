import Infos from '@/pages/Labels/infos/infos';
import { getLabelDetail } from '@/services/api';
import { useParams } from '@@/exports';
import ProCard from '@ant-design/pro-card';
import { PageContainer, ProDescriptions } from '@ant-design/pro-components';
import { message } from 'antd';
import { useEffect, useState } from 'react';

export default function () {
  const params = useParams();

  const [detail, setDetail] = useState({
    name: '',
    code: '',
    description: '',
    updated_at: '',
    created_at: '',
  });

  const getData = async () => {
    const result = await getLabelDetail(params?.id);
    if (result) {
      setDetail(result.data);
    } else {
      message.error('获取详细信息失败');
    }
  };

  useEffect(() => {
    getData();
  }, []);
  console.log(params);
  return (
    <PageContainer
      header={{
        title: '标签信息',
      }}
    >
      <ProCard>
        <ProDescriptions title={'基础信息'} column={4}>
          <ProDescriptions.Item label="名称" span={1}>
            {detail?.name}
          </ProDescriptions.Item>
          <ProDescriptions.Item label="代码" span={1}>
            {detail?.code}
          </ProDescriptions.Item>

          <ProDescriptions.Item label="创建时间" span={1} valueType="dateTime">
            {detail?.created_at}
          </ProDescriptions.Item>
          <ProDescriptions.Item label="更新时间" span={1} valueType="dateTime">
            {detail?.updated_at}
          </ProDescriptions.Item>
          <ProDescriptions.Item label="描述" span={4}>
            {detail?.description}
          </ProDescriptions.Item>
        </ProDescriptions>
      </ProCard>
      <Infos labelID={params?.id} />
    </PageContainer>
  );
}
