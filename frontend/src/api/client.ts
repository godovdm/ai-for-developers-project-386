import { DefaultApi } from "./generated/apis/DefaultApi";
import { Configuration, ResponseError } from "./generated/runtime";

export const api = new DefaultApi(new Configuration({ basePath: "" }));

export async function readErrorMessage(
  error: unknown,
  fallback: string,
): Promise<string> {
  if (error instanceof ResponseError) {
    try {
      const body: unknown = await error.response.clone().json();
      if (
        typeof body === "object" &&
        body !== null &&
        "message" in body &&
        typeof body.message === "string"
      ) {
        return body.message;
      }
    } catch {
      // The response body is not JSON; use the fallback message.
    }
  }
  return fallback;
}
