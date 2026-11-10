import { Anchor, type AnchorProps } from "@mantine/core";
import { useEffect, useState, type MouseEvent, type ReactNode } from "react";

const NAVIGATE_EVENT = "calendar:navigate";

export function usePath(): string {
  const [path, setPath] = useState(() => window.location.pathname);

  useEffect(() => {
    const update = () => {
      setPath(window.location.pathname);
    };
    window.addEventListener("popstate", update);
    window.addEventListener(NAVIGATE_EVENT, update);
    return () => {
      window.removeEventListener("popstate", update);
      window.removeEventListener(NAVIGATE_EVENT, update);
    };
  }, []);

  return path;
}

export function navigate(to: string): void {
  if (window.location.pathname === to) return;
  window.history.pushState({}, "", to);
  window.dispatchEvent(new Event(NAVIGATE_EVENT));
  window.scrollTo({ top: 0 });
}

type AppLinkProps = {
  to: string;
  children: ReactNode;
} & Omit<AnchorProps, "href" | "onClick" | "children">;

export function AppLink({ to, children, ...props }: AppLinkProps) {
  const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
    if (
      event.defaultPrevented ||
      event.button !== 0 ||
      event.metaKey ||
      event.ctrlKey ||
      event.shiftKey ||
      event.altKey
    ) {
      return;
    }
    event.preventDefault();
    navigate(to);
  };

  return (
    <Anchor href={to} onClick={handleClick} {...props}>
      {children}
    </Anchor>
  );
}
