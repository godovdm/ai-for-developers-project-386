import { act, fireEvent, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, expect, it, vi } from "vitest";
import { MeetingDuration } from "./api/generated/models/MeetingDuration";
import { renderApp, jsonResponse } from "./test/utils";

beforeEach(() => {
  window.history.replaceState({}, "", "/");
});

afterEach(() => {
  vi.unstubAllGlobals();
});

it("shows the guest page with meeting types after loading", async () => {
  let resolveFetch: ((response: Response) => void) | undefined;
  const fetchMock = vi.fn<typeof fetch>().mockImplementation(
    () =>
      new Promise<Response>((resolve) => {
        resolveFetch = resolve;
      }),
  );
  vi.stubGlobal("fetch", fetchMock);

  renderApp("/");

  expect(
    screen.getByRole("heading", { name: "Выберите время для встречи" }),
  ).toBeInTheDocument();
  expect(screen.getByLabelText("Загрузка типов встреч")).toBeInTheDocument();

  await waitFor(() => {
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });

  await act(async () => {
    resolveFetch?.(
      jsonResponse({
        items: [
          {
            id: "1",
            name: "Вводный звонок",
            description: "Знакомство с командой",
            durationMinutes: MeetingDuration.NUMBER_60,
          },
        ],
      }),
    );
    await Promise.resolve();
  });

  expect(
    await screen.findByRole("heading", { name: "Вводный звонок" }),
  ).toBeInTheDocument();
  expect(screen.getByText("Знакомство с командой")).toBeInTheDocument();
  expect(screen.getByText("60 минут")).toBeInTheDocument();
  expect(
    screen.getByRole("button", { name: "Выбрать время" }),
  ).toBeInTheDocument();
});

it("shows the empty state and keeps the owner link", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn<typeof fetch>().mockResolvedValue(jsonResponse({ items: [] })),
  );

  renderApp("/");

  expect(
    await screen.findByText("Пока нет доступных встреч"),
  ).toBeInTheDocument();
  expect(screen.getByRole("link", { name: "Владельцу" })).toBeInTheDocument();
});

it("shows the not found page for a malformed meeting type URL", () => {
  window.history.pushState({}, "", "/meeting-types/%FF");

  renderApp("/meeting-types/%FF");

  expect(
    screen.getByRole("heading", { name: "Страница не найдена" }),
  ).toBeInTheDocument();
  expect(screen.getByRole("link", { name: "Calendar" })).toBeInTheDocument();
});

it("shows an error and retries loading meeting types", async () => {
  const fetchMock = vi
    .fn<typeof fetch>()
    .mockRejectedValueOnce(new TypeError("offline"))
    .mockResolvedValueOnce(
      jsonResponse({
        items: [
          {
            id: "1",
            name: "Вводный звонок",
            description: "Знакомство с командой",
            durationMinutes: MeetingDuration.NUMBER_30,
          },
        ],
      }),
    );
  vi.stubGlobal("fetch", fetchMock);

  renderApp("/");

  expect(
    await screen.findByText("Не удалось загрузить типы встреч"),
  ).toBeInTheDocument();

  fireEvent.click(screen.getByRole("button", { name: "Попробовать снова" }));

  expect(
    await screen.findByRole("heading", { name: "Вводный звонок" }),
  ).toBeInTheDocument();
  expect(fetchMock).toHaveBeenCalledTimes(2);
});
