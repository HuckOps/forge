import { message } from 'antd';
import { extend } from 'umi-request';

const errorHandler = async (error: any) => {
  const { response } = error;
  console.log(response);
  if (response && response.status) {
    if (response.status < 400) {
      // For successful responses (<=400), return the response data
      console.log(response);
      return response.data;
    } else {
      // For error responses (>400), show error message and return false
      try {
        const errorData = await response.clone().json();
        const errorMsg =
          errorData.msg || errorData.message || response.statusText;
        message.error(errorMsg);
      } catch (e) {
        message.error(response.statusText);
      }
      return false;
    }
  } else if (!response) {
    message.error('Network error');
    return false;
  }

  return response;
};

export const request = extend({
  errorHandler,
  // Parse response JSON automatically
  // parseResponse: true,
});
