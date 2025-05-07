import { DynamoDBClient } from "@aws-sdk/client-dynamodb";
import { CustomEvent, CustomResponse, RequestBody } from "./types";

const dynamoDb = new DynamoDBClient({});

export const handler = async (event: CustomEvent): Promise<CustomResponse> => {
  try {
    // リクエストの解析
    const { httpMethod, path, body, queryStringParameters } = event;

    // リクエストボディの解析（存在する場合）
    let parsedBody: RequestBody | undefined;
    if (body) {
      try {
        parsedBody = JSON.parse(body) as RequestBody;
      } catch (e) {
        return {
          statusCode: 400,
          headers: {
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "*",
          },
          body: JSON.stringify({
            message: "Invalid request body",
            error: "Failed to parse request body",
          }),
        };
      }
    }

    // レスポンスの基本構造
    const response: CustomResponse = {
      statusCode: 200,
      headers: {
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin": "*", // CORS対応
      },
      body: JSON.stringify({
        message: "Hello from Lambda!",
        path,
        method: httpMethod,
        timestamp: new Date().toISOString(),
        requestData: {
          body: parsedBody,
          query: queryStringParameters,
        },
        requestId: event.requestContext.requestId,
      }),
    };

    return response;
  } catch (error) {
    console.error("Error:", error);
    return {
      statusCode: 500,
      headers: {
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin": "*",
      },
      body: JSON.stringify({
        message: "Internal server error",
        error: error instanceof Error ? error.message : "Unknown error",
      }),
    };
  }
};
