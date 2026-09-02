import { type } from "arktype";

export const Ding_DTO_Schema = type({
  id: "number",
  name: "string > 0",
  angelegt: "string.date.iso.parse",
});
export type Ding_DTO = typeof Ding_DTO_Schema.infer;

export const Nur_Request_DTO_Schema = type({
  wann: "Date",
});
export type Nur_Request_DTO = typeof Nur_Request_DTO_Schema.infer;

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
  zeitpunkt: "string.date.iso.parse",
  vielleicht: "string.date.iso.parse?",
  oderNull: "string.date.iso.parse | null",
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
  nurRequest: Nur_Request_DTO_Schema,
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
      log?: { warn(meldung: unknown, ...details: unknown[]): void };
    },
  ) { }

  async #call<TRequest, TResponse>(
    path: string,
    args: TRequest,
    schema: (data: unknown) => unknown,
  ): Promise<{ value: TResponse; error: null; status: number } | { value: null; error: string; status: number | null; body: unknown }> {

    const do_fetch = this.options?.fetch ?? globalThis.fetch;
    const do_log = this.options?.log ?? console;

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
        do_log.warn(`Fetch error: ${result.status} ${result.statusText} for ${path}`);
        if (this.options?.handle_error) this.options.handle_error(result);
        const fehler_body = await result.json().catch(() => null);
        return {
          value: null,
          error: (fehler_body as { message?: string } | null)?.message ?? 'Unknown error',
          status: result.status,
          body: fehler_body,
        };
      }

      const geprueft = schema(await result.json());

      if (geprueft instanceof type.errors) {
        const stellen = geprueft.map((fehler) => `${fehler.path.join(".")} erwartet ${fehler.expected}`).join("; ");
        do_log.warn(`Antwort verletzt Vertrag: ${path} — ${stellen}`);

        return {
          value: null,
          error: "Antwort verletzt den Vertrag",
          status: null,
          body: null,
        };
      }

      return {
        value: geprueft as TResponse,
        error: null,
        status: result.status,
      };
    } catch (error) {
      do_log.warn(`RPC_Client Error for ${path}`);
      do_log.warn(error);

      return {
        value: null,
        error: error instanceof Error ? error.message : "Unknown error",
        status: null,
        body: null,
      };
    }
  }

  a_name = (args: A_Name_Request) =>
    this.#call<A_Name_Request, A_Name_Response>(A_Name_Path, args, A_Name_Response_Schema);

  eins = (args: Eins_Request) =>
    this.#call<Eins_Request, Eins_Response>(Eins_Path, args, Eins_Response_Schema);

  listen = (args: Listen_Request) =>
    this.#call<Listen_Request, Listen_Response>(Listen_Path, args, Listen_Response_Schema);

  zwei = (args: Zwei_Request) =>
    this.#call<Zwei_Request, Zwei_Response>(Zwei_Path, args, Zwei_Response_Schema);
}
