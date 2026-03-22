type NoVncDisconnectEvent = {
  detail?: {
    clean?: boolean;
  };
};

type NoVncCredentials = {
  password: string;
};

export interface NoVncRfb {
  scaleViewport: boolean;
  resizeSession: boolean;
  clipViewport: boolean;
  background: string;
  addEventListener(
    event: "connect" | "credentialsrequired",
    listener: () => void,
  ): void;
  addEventListener(
    event: "disconnect",
    listener: (event: NoVncDisconnectEvent) => void,
  ): void;
  sendCredentials(credentials: NoVncCredentials): void;
  disconnect(): void;
}

export type NoVncConstructor = new (
  target: Element,
  url: string,
  options?: { credentials?: NoVncCredentials },
) => NoVncRfb;

export type NoVncModule = {
  default: NoVncConstructor;
};

declare module "/vendor/novnc/rfb.bundle.mjs" {
  const novncBundle: NoVncModule;
  export default novncBundle.default;
}
