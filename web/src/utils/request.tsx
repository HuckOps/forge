import { message } from 'antd';
import { extend } from 'umi-request';
const errorHandler = (error: any) => {
  const { response } = error;
  if (response && response.status) {
    const errorText = response.statusText;
    const { status, url } = response;

    message.error(status, url, errorText);
    return Promise.reject(false);
  } else if (!response) {
    message.error('network error');
    return Promise.reject(false);
  }
  return response;
};
export const request = extend(errorHandler);
