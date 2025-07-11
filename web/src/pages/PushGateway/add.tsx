import {
  createPushGateway,
  getNodes,
  getPushGateway,
  updatePushGateway,
} from '@/services/api';
import ProCard from '@ant-design/pro-card';
import { PageContainer, ProFormInstance } from '@ant-design/pro-components';
import {
  ProForm,
  ProFormDigit,
  ProFormSelect,
  ProFormText,
} from '@ant-design/pro-form/lib';
import { history, useLocation } from '@umijs/max';
import { message } from 'antd';
import { useEffect, useRef, useState } from 'react';

export default function () {
  const { search } = useLocation();
  let searchParams = new URLSearchParams(search);
  const id = searchParams.get('id');
  const ref = useRef<ProFormInstance>(null);
  const [nodes, setNodes] = useState([]);

  const getNodesByIP = async (ip: string) => {
    const result = await getNodes(0, -1, {
      ip: ip,
    });
    if (result) {
      setNodes(
        result.data.data?.map((item: { ip: any; uuid: any }) => ({
          label: item.ip,
          value: item.uuid,
        })),
      );
    }
  };

  const onSubmit = async () => {
    const data = await ref.current?.validateFields();
    var result;
    if (id) {
      result = await updatePushGateway(id, data);
    } else {
      result = await createPushGateway(data);
    }
    console.log(result);
    if (result) {
      message.success('操作PushGateway成功');
      history.push('/prometheus/pushgateway');
    } else {
      message.error('操作PushGateway失败');
    }
  };

  useEffect(() => {
    if (id) {
      getPushGateway(id).then((res) => {
        console.log(res?.data);
        if (res?.data) {
          ref?.current?.setFieldsValue(res?.data);
          getNodesByIP(res?.data?.ip);
        }
      });
    }
  }, []);

  return (
    <PageContainer>
      <ProCard>
        <ProForm
          layout="horizontal"
          formRef={ref}
          submitter={{
            onSubmit: onSubmit,
          }}
        >
          <ProFormDigit
            label="端口"
            name="port"
            rules={[
              {
                required: true,
                message: '端口为必填项',
              },
            ]}
          />
          <ProFormText
            label="版本"
            name="version"
            rules={[
              {
                required: true,
                message: '版本为必填项',
              },
            ]}
          />
          <ProFormSelect
            label="节点"
            name="uuid"
            disabled={!!id}
            showSearch={true}
            fieldProps={{
              placeholder: '选择目标部署机器，仅支持ip搜索',
              onSearch: async (value: string) => {
                if (value) await getNodesByIP(value);
                else message.error('查询机器失败');
              },
            }}
            options={nodes}
            rules={[
              {
                required: true,
                message: '机器为必选项',
              },
            ]}
          />
        </ProForm>
      </ProCard>
    </PageContainer>
  );
}
