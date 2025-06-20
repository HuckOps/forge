import { request } from '@/utils/request';

export const getNodes = async (
  skip: number = 0,
  limit: number = 10,
  filter: any = {},
) => {
  console.log(filter);
  return await request.get(`/api/v1/nodes`, {
    params: {
      skip,
      limit,
      ...filter,
    },
  });
};
