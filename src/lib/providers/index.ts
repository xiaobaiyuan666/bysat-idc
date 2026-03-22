import { MofangCloudProvider } from "@/lib/providers/mofang-cloud";

const provider = new MofangCloudProvider();

export function getCloudProvider() {
  return provider;
}
