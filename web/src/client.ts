import { createConnectTransport } from "@connectrpc/connect-web";
import { createPromiseClient } from "@connectrpc/connect";
import { TaxService } from "./gen/api/tax_connectweb";

const baseUrl = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

const transport = createConnectTransport({ baseUrl });
export const taxClient = createPromiseClient(TaxService, transport);