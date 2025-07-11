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

export const getLabels = async (
  skip: number = 0,
  limit: number = 10,
  filter: any = {},
): Promise<API.Restful<API.Pagination<API.Label>>> => {
  return await request.get(`/api/v1/labels`, {
    params: {
      skip,
      limit,
      ...filter,
    },
  });
};

export const createLabel = async (data: any = {}) => {
  return await request.post(`/api/v1/labels`, {
    data,
  });
};

export const getLabelDetail = async (id: any) => {
  return await request.get(`/api/v1/labels/${id}`);
};

export const setNodeLabel = async (data: any) => {
  return await request.post(`/api/v1/nodes/labels`, {
    data,
  });
};

export const getLabelNodes = async (
  id?: string,
  skip: number = 0,
  limit: number = 10,
) => {
  return await request.get(`/api/v1/labels/${id}/nodes`, {
    params: {
      skip,
      limit,
    },
  });
};

export const getPushGateways = async (skip: number, limit: number) => {
  return await request.get(`/api/v1/prometheus/pushgateway`, {
    params: {
      skip,
      limit,
    },
  });
};

export const createPushGateway = async (data: any = {}) => {
  return await request.post(`/api/v1/prometheus/pushgateway`, {
    data,
  });
};

export const deletePushGateway = async (id: any) => {
  return await request.delete(`/api/v1/prometheus/pushgateway/${id}`);
};

export const getPushGateway = async (id: any) => {
  return await request.get(`/api/v1/prometheus/pushgateway/${id}`);
};

export const updatePushGateway = async (id: any, data: any) => {
  return await request.put(`/api/v1/prometheus/pushgateway/${id}`, {
    data,
  });
};

export const createFederation = async (data: any) => {
  return await request.post(`/api/v1/prometheus/federation`, {
    data,
  });
};

export const updateFederation = async (id: any, data: any) => {
  return await request.put(`/api/v1/prometheus/federation/${id}`, {
    data,
  });
};

export const getFederations = async (skip = 0, limit = 10) => {
  return await request.get(`/api/v1/prometheus/federation`, {
    params: {
      skip,
      limit,
    },
  });
};
