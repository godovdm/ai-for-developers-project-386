import { MantineProvider } from "@mantine/core";
import { render, screen } from "@testing-library/react";
import { afterEach, beforeEach, expect, it, vi } from "vitest";
import { App } from "./App";

beforeEach(() => {
  window.matchMedia = vi.fn().mockReturnValue({
    matches: false,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
  });
});

afterEach(() => {
  vi.unstubAllGlobals();
});

it("shows that the server is running", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ status: "ok" }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      }),
    ),
  );
  render(
    <MantineProvider>
      <App />
    </MantineProvider>,
  );
  expect(await screen.findByText("Сервер работает")).toBeTruthy();
});

it("shows an error when the server is unavailable", async () => {
  vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new Error("offline")));
  render(
    <MantineProvider>
      <App />
    </MantineProvider>,
  );
  expect(await screen.findByText("Сервер недоступен")).toBeTruthy();
});
