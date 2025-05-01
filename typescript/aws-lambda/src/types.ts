// リクエストボディの型定義
export interface RequestBody {
  userId?: string;
  data?: Record<string, unknown>;
}

// カスタムイベント型の定義
export interface CustomEvent {
  body: string | null;
  headers: Record<string, string | undefined>;
  httpMethod: string;
  path: string;
  queryStringParameters: Record<string, string | undefined> | null;
  requestContext: {
    requestId: string;
    timeEpoch: number;
  };
}

// レスポンス型の定義
export interface CustomResponse {
  statusCode: number;
  headers: Record<string, string>;
  body: string;
}
