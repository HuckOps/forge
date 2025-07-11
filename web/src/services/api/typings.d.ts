declare namespace API {
  interface Restful<T> {
    code: int;
    msg: string;
    data: T;
  }
  interface Pagination<T> {
    total: number;
    data: T[];
  }

  interface Node {
    id: string;
    created_at: string;
    updated_at: string;
    uuid: string;
    hostname: string;
    heartbeat: string;
    heartbeat_status: boolean;
  }

  interface Label {
    id: string;
    name: string;
    code: string;
    created_at: string;
    updated_at: string;
  }
}
