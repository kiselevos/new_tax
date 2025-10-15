import { createConnectTransport } from "@connectrpc/connect-web";
import { createPromiseClient } from "@connectrpc/connect";
import { TaxService } from "./gen/api/tax_connect";

const baseUrl = import.meta.env.VITE_API_URL || "/";

const transport = createConnectTransport({
  baseUrl,
  useBinaryFormat: false,
});

export const taxClient = createPromiseClient(TaxService, transport);