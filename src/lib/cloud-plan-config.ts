export type CloudPlanConfig = {
  node?: string;
  os?: string;
  cpu?: number;
  memory?: number;
  system_disk_size?: number;
  data_disk_size?: number;
  network_type?: string;
  bw?: number;
  in_bw?: number;
  flow_limit?: number;
  flow_way?: string;
  ip_num?: number;
  peak_defence?: number;
  [key: string]: unknown;
};

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
}

function normalizePositiveInt(value?: number | null) {
  if (value === undefined || value === null) {
    return undefined;
  }

  if (!Number.isFinite(value)) {
    return undefined;
  }

  const normalized = Math.round(value);
  return normalized > 0 ? normalized : undefined;
}

export function parseCloudPlanConfig(input?: string | null): CloudPlanConfig {
  if (!input) {
    return {};
  }

  try {
    const parsed = JSON.parse(input) as CloudPlanConfig;
    return parsed && typeof parsed === "object" ? parsed : {};
  } catch {
    return {};
  }
}

export function buildCloudPlanConfig(input: {
  providerNodeId?: string | null;
  providerOsCode?: string | null;
  configCpu?: number | null;
  configMemoryGb?: number | null;
  systemDiskSize?: number | null;
  dataDiskSize?: number | null;
  networkType?: string | null;
  bandwidthMbps?: number | null;
  inboundBandwidthMbps?: number | null;
  flowLimitGb?: number | null;
  flowBillingMode?: string | null;
  ipCount?: number | null;
  peakDefenceGbps?: number | null;
  extra?: Record<string, unknown>;
}) {
  const config: CloudPlanConfig = {
    ...(input.extra ?? {}),
  };

  const node = normalizeText(input.providerNodeId);
  const os = normalizeText(input.providerOsCode);
  const networkType = normalizeText(input.networkType);
  const flowWay = normalizeText(input.flowBillingMode);
  const cpu = normalizePositiveInt(input.configCpu);
  const memory = normalizePositiveInt(input.configMemoryGb);
  const systemDiskSize = normalizePositiveInt(input.systemDiskSize);
  const dataDiskSize = normalizePositiveInt(input.dataDiskSize);
  const bandwidth = normalizePositiveInt(input.bandwidthMbps);
  const inboundBandwidth = normalizePositiveInt(input.inboundBandwidthMbps);
  const flowLimit = normalizePositiveInt(input.flowLimitGb);
  const ipCount = normalizePositiveInt(input.ipCount);
  const peakDefence = normalizePositiveInt(input.peakDefenceGbps);

  if (node) {
    config.node = node;
  }

  if (os) {
    config.os = os;
  }

  if (cpu) {
    config.cpu = cpu;
  }

  if (memory) {
    config.memory = memory;
  }

  if (systemDiskSize) {
    config.system_disk_size = systemDiskSize;
  }

  if (dataDiskSize) {
    config.data_disk_size = dataDiskSize;
  }

  if (networkType) {
    config.network_type = networkType;
  }

  if (bandwidth) {
    config.bw = bandwidth;
  }

  if (inboundBandwidth) {
    config.in_bw = inboundBandwidth;
  }

  if (flowLimit) {
    config.flow_limit = flowLimit;
  }

  if (flowWay) {
    config.flow_way = flowWay;
  }

  if (ipCount) {
    config.ip_num = ipCount;
  }

  if (peakDefence) {
    config.peak_defence = peakDefence;
  }

  return Object.keys(config).length > 0 ? JSON.stringify(config) : undefined;
}
