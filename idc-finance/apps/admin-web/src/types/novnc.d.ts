export interface NoVncRfb {
  scaleViewport: boolean;
  resizeSession: boolean;
  clipViewport: boolean;
  background: string;
  addEventListener(type: string, listener: (event: any) => void): void;
  disconnect(): void;
  sendCredentials(credentials: { password?: string }): void;
}

export interface NoVncModule {
  default: new (
    target: Element,
    url: string,
    options?: { credentials?: { password?: string } }
  ) => NoVncRfb;
}
