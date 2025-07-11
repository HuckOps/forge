import { createLabel } from '@/services/api';
import ProCard from '@ant-design/pro-card';
import {
  PageContainer,
  ProFormInstance,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { ProForm, ProFormText } from '@ant-design/pro-form/lib';
import { history } from '@umijs/max';
import { message } from 'antd';
import { useRef } from 'react';

export default function () {
  const ref = useRef<ProFormInstance>();
  const onSubmit = async () => {
    const data = await ref.current?.validateFields();
    const result = await createLabel(data);
    console.log(result);
    if (result) {
      message.info('创建标签成功');
      history.push('/nodes/labels');
    } else {
      message.error('创建标签失败');
    }
  };
  return (
    <PageContainer
      header={{
        title: '添加标签',
      }}
    >
      <ProCard>
        <ProForm
          formRef={ref}
          layout="horizontal"
          grid
          submitter={{
            onSubmit: onSubmit,
            searchConfig: {
              submitText: '提交',
            },
            resetButtonProps: {
              style: {
                display: 'none',
              },
            },
            render: (_, dom) => (
              <div style={{ textAlign: 'center' }}>{dom}</div>
            ),
          }}
        >
          <ProFormText
            label="名称"
            name="name"
            colProps={{ md: 12, xl: 12 }}
            rules={[{ required: true, message: '请输入名称' }]}
            initialValue={''}
          />
          <ProFormText
            label="代号"
            name="code"
            colProps={{ md: 12, xl: 12 }}
            rules={[
              { required: true, message: '请输入代号' },
              {
                pattern: /^[a-zA-Z0-9_]+$/,
                message: '代号只能包含字母、数字和下划线',
              },
            ]}
            initialValue={''}
          />
          <ProFormTextArea
            label="描述"
            name="description"
            width={'xl'}
            colProps={{
              md: 24,
            }}
            initialValue={''}
          />
        </ProForm>
      </ProCard>
    </PageContainer>
  );
}
