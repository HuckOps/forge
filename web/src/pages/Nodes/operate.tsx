import { getLabels, setNodeLabel } from '@/services/api';
import { ProFormInstance } from '@ant-design/pro-components';
import { ProForm, ProFormRadio, ProFormSelect } from '@ant-design/pro-form/lib';
import { message } from 'antd';
import { useEffect, useRef, useState } from 'react';

interface Props {
  SelectedRowKeys: string[];
  resetSelectRowKeys: () => void;
}

export default function (props: Props) {
  const ref = useRef<ProFormInstance>(null);
  const [operate, setOperate] = useState<string>();
  const [labels, setLabels] = useState<API.Label[]>([]);
  useEffect(() => {
    getLabels(0, -1).then((res) => {
      return setLabels(res.data.data);
    });
  }, []);

  const onSubmit = async () => {
    const data = await ref?.current?.validateFields();
    switch (data.operate) {
      case 'add_labels': {
        const result = await setNodeLabel({
          nodes: props.SelectedRowKeys,
          labels: data.labels,
        });
        if (result) {
          message.success('配置成功');
          props.resetSelectRowKeys();
        }
      }
    }
  };

  if (props.SelectedRowKeys.length !== 0) {
    return (
      <ProForm
        layout="horizontal"
        formRef={ref}
        onValuesChange={(changedValues) => {
          if ('operate' in changedValues) {
            setOperate(changedValues.operate);
          }
        }}
        submitter={{
          onSubmit: onSubmit,
        }}
      >
        <ProFormRadio.Group
          label="选择操作"
          name="operate"
          options={[
            { label: '添加标签', value: 'add_labels' },
            {
              label: '删除标签',
              value: 'del_labels',
            },
          ]}
          rules={[{ required: true, message: '请选择操作' }]}
        ></ProFormRadio.Group>
        {operate === 'add_labels' ? (
          <ProFormSelect
            label="选择标签"
            name="labels"
            mode="multiple"
            options={labels?.map((label) => ({
              label: label?.name,
              value: label?.id,
            }))}
            rules={[{ required: true, message: '请选择标签' }]}
          ></ProFormSelect>
        ) : null}
      </ProForm>
    );
  }
}
