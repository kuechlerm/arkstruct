import { type } from "arktype";

export const Ding_DTO_Schema = type({
  id: "number",
  name: "string > 0",
});
export type Ding_DTO = typeof Ding_DTO_Schema.infer;

export const A_Name_Path = "/a_name";
export const A_Name_Request_Schema = type({
  msg: "string > 0",
});
export type A_Name_Request = typeof A_Name_Request_Schema.infer;

export const A_Name_Response_Schema = type({
  msg: "string > 0",
});
export type A_Name_Response = typeof A_Name_Response_Schema.infer;

export const Eins_Path = "/eins";
export const Eins_Request_Schema = type({
  requiredString: "string > 0",
  optionalString: "string | undefined",
  requiredInt: "number > 0",
  optionalInt: "number | undefined",
  requiredBool: "boolean",
  optionalBool: "boolean | undefined",
  specialType: type.or(A_Name_Request_Schema, A_Name_Response_Schema),
});
export type Eins_Request = typeof Eins_Request_Schema.infer;

export const Eins_Response_Schema = type({
  responseString: "string > 0",
});
export type Eins_Response = typeof Eins_Response_Schema.infer;

export const Listen_Path = "/listen";
export const Listen_Request_Schema = type({});
export type Listen_Request = typeof Listen_Request_Schema.infer;

export const Listen_Response_Schema = type({
  dinge: Ding_DTO_Schema.array(),
});
export type Listen_Response = typeof Listen_Response_Schema.infer;

export const Zwei_Path = "/zwei";
export const Zwei_Request_Schema = type({
  optionalString: "string | undefined",
});
export type Zwei_Request = typeof Zwei_Request_Schema.infer;

export const Zwei_Response_Schema = type({
  responseString: "string > 0",
});
export type Zwei_Response = typeof Zwei_Response_Schema.infer;

export class RPC_Client {
  constructor(
    private base_url: string,
    private options?: {
      fetch?: (url: string, init: RequestInit) => Promise<Response>;
      handle_error?: (response: Response) => void;
    },
  ) { }

  async #call<TRequest, TResponse>(
    path: string,
    args: TRequest,
  ): Promise<{ value: TResponse; error: null; status: number } | { value: null; error: string; status: number | null }> {

    const do_fetch = this.options?.fetch ?? globalThis.fetch;

    try {
      const result = await do_fetch(new URL(path, this.base_url).href, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(args),
      });

      if (!result.ok) {
        console.warn(`Fetch error: ${result.status} ${result.statusText} for ${path}`);
        if (this.options?.handle_error) this.options.handle_error(result);
        return {
          value: null,
          error: (await result.json())?.message ?? 'Unknown error',
          status: result.status,
        };
      }

      const data = await result.json();
      const revived = this.revive_dates(data);

      return {
        value: revived as TResponse,
        error: null,
        status: result.status,
      };
    } catch (error) {
      console.warn('RPC_Client Error for', { path, args: JSON.stringify(args) });
      console.warn(error);

      return {
        value: null,
        error: error instanceof Error ? error.message : "Unknown error",
        status: null,
      };
    }
  }

  revive_dates = <T>(obj: T): T => {
    const ISO_DATE_REGEX = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?([+-]\d{2}:\d{2}|Z)$/;

    if (obj == null || typeof obj !== 'object') return obj;

    if (Array.isArray(obj)) {
      return obj.map(this.revive_dates) as any;
    }

    const result: any = {};
    for (const [key, value] of Object.entries(obj)) {
      if (typeof value === 'string' && ISO_DATE_REGEX.test(value)) {
        result[key] = new Date(value);
      } else if (typeof value === 'object' && value !== null) {
        result[key] = this.revive_dates(value);
      } else {
        result[key] = value;
      }
    }
    return result;
  };

  a_name = (args: A_Name_Request) =>
    this.#call<A_Name_Request, A_Name_Response>(A_Name_Path, args);

  eins = (args: Eins_Request) =>
    this.#call<Eins_Request, Eins_Response>(Eins_Path, args);

  listen = (args: Listen_Request) =>
    this.#call<Listen_Request, Listen_Response>(Listen_Path, args);

  zwei = (args: Zwei_Request) =>
    this.#call<Zwei_Request, Zwei_Response>(Zwei_Path, args);
}
