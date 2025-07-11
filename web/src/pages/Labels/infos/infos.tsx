import Nodes from '@/pages/Labels/infos/nodes';
import ProCard from '@ant-design/pro-card';
import { Tabs } from 'antd';
import { useEffect, useState } from 'react';

type Props = {
  labelID?: string;
};

export default function (props: Props) {
  const [tabsActiveKey, setTabsActiveKey] = useState<string>('nodes');
  useEffect(() => {
    switch (tabsActiveKey) {
      case 'nodes': {
      }
    }
  }, [tabsActiveKey]);
  const renderInfos = () => {
    switch (tabsActiveKey) {
      case 'nodes': {
        return <Nodes labelID={props.labelID} />;
      }
    }
  };
  return (
    <>
      <ProCard>
        <Tabs
          activeKey={tabsActiveKey}
          onChange={(activeKey) => setTabsActiveKey(activeKey)}
          items={[
            {
              label: '节点列表',
              key: 'nodes',
            },
          ]}
        />
        {renderInfos()}
      </ProCard>
    </>
  );
}
