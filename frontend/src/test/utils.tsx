import { MantineProvider } from "@mantine/core";
import { render } from "@testing-library/react";
import { App } from "../App";

export function renderApp(path: string) {
  window.history.replaceState({}, "", path);
  return render(
    <MantineProvider defaultColorScheme="light">
      <App />
    </MantineProvider>,
  );
}

export function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}
