export type ProviderAction =
  | "provision"
  | "sync"
  | "activate"
  | "suspend"
  | "unsuspend"
  | "renew"
  | "terminate"
  | "powerOn"
  | "powerOff"
  | "reboot"
  | "hardReboot"
  | "hardPowerOff"
  | "reinstall"
  | "resetPassword"
  | "getVnc"
  | "rescueStart"
  | "rescueStop"
  | "lock"
  | "unlock";

export type ProviderActionPayload = {
  imageId?: string;
  password?: string;
  notes?: string;
  rescueType?: string;
  [key: string]: unknown;
};

export type ProviderServicePayload = {
  localId: string;
  serviceNo: string;
  name: string;
  hostname?: string | null;
  providerResourceId?: string | null;
  providerProductId?: string | null;
  region?: string | null;
  billingCycle: string;
  customerName: string;
  productName: string;
  configOptions?: Record<string, unknown>;
};

export type ProviderResponse = {
  ok: boolean;
  status: string;
  remoteId?: string;
  ipAddress?: string;
  region?: string;
  nextDueDate?: Date;
  message?: string;
  requestBody?: string;
  responseBody?: string;
  consoleUrl?: string;
  taskId?: string;
  data?: Record<string, unknown>;
};

export interface CloudProvider {
  name: string;
  listInstances(): Promise<ProviderResponse[]>;
  provisionService(service: ProviderServicePayload): Promise<ProviderResponse>;
  activateService(service: ProviderServicePayload): Promise<ProviderResponse>;
  suspendService(service: ProviderServicePayload): Promise<ProviderResponse>;
  renewService(service: ProviderServicePayload): Promise<ProviderResponse>;
  terminateService(service: ProviderServicePayload): Promise<ProviderResponse>;
  syncService(service: ProviderServicePayload): Promise<ProviderResponse>;
  manageServiceAction(
    service: ProviderServicePayload,
    action: ProviderAction,
    payload?: ProviderActionPayload,
  ): Promise<ProviderResponse>;
}
