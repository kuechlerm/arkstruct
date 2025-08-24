import { type } from "arktype";

export const Ding_DTO_Schema = type({
  id: "number",
  name: "string > 0",
});
export type Ding_DTO = typeof Ding_DTO_Schema.infer;

export const A_Name_Request_Schema = type({
  msg: "string > 0",
});
export type A_Name_Request = typeof A_Name_Request_Schema.infer;

export const A_Name_Response_Schema = type({
  msg: "string > 0",
});
export type A_Name_Response = typeof A_Name_Response_Schema.infer;

export const Eins_Request_Schema = type({
  requiredString: "string > 0",
  optionalString: "string | undefined",
  requiredInt: "number > 0",
  optionalInt: "number | undefined",
  requiredBool: "boolean",
  optionalBool: "boolean | undefined",
});
export type Eins_Request = typeof Eins_Request_Schema.infer;

export const Eins_Response_Schema = type({
  responseString: "string > 0",
});
export type Eins_Response = typeof Eins_Response_Schema.infer;

export const Listen_Request_Schema = type({});
export type Listen_Request = typeof Listen_Request_Schema.infer;

export const Listen_Response_Schema = type({
  dinge: Ding_DTO_Schema.array(),
});
export type Listen_Response = typeof Listen_Response_Schema.infer;

export const Zwei_Request_Schema = type({
  optionalString: "string | undefined",
});
export type Zwei_Request = typeof Zwei_Request_Schema.infer;

export const Zwei_Response_Schema = type({
  responseString: "string > 0",
});
export type Zwei_Response = typeof Zwei_Response_Schema.infer;

export class RPC_Client {
  constructor(public base_url: string) {}

  async #call<TRequest, TResponse>(
    path: string,
    args: TRequest,
  ): Promise<
    { value: TResponse; error: null } | { value: null; error: string }
  > {
    try {
      const result = await fetch(new URL(path, this.base_url).href, {
        method: "POST",
        body: JSON.stringify(args),
      });

      if (!result.ok) {
        console.error(
          `Fetch error: ${result.status} ${result.statusText} for ${path}`,
        );
        return {
          value: null,
          error: `Fetch error: ${result.status} ${result.statusText}`,
        };
      }

      const data = await result.json();

      return {
        value: data as TResponse,
        error: null,
      };
    } catch (error) {
      console.error(`Error during fetch for ${path}:`, error);

      return {
        value: null,
        error: error instanceof Error ? error.message : "Unknown error",
      };
    }
  }

  a_name = (args: A_Name_Request) =>
    this.#call<A_Name_Request, A_Name_Response>("/a_name", args);

  eins = (args: Eins_Request) =>
    this.#call<Eins_Request, Eins_Response>("/eins", args);

  listen = (args: Listen_Request) =>
    this.#call<Listen_Request, Listen_Response>("/listen", args);

  zwei = (args: Zwei_Request) =>
    this.#call<Zwei_Request, Zwei_Response>("/zwei", args);
}
