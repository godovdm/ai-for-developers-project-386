import { fireEvent, screen } from "@testing-library/react";
import { afterEach, beforeEach, expect, it, vi } from "vitest";
import { MeetingDuration } from "../api/generated/models/MeetingDuration";
import type { MeetingType } from "../api/generated/models/MeetingType";
import { jsonResponse, renderApp } from "../test/utils";

beforeEach(() => {
  window.history.replaceState({}, "", "/owner");
});

afterEach(() => {
  vi.unstubAllGlobals();
});

function requestUrl(input: RequestInfo | URL): string {
  if (typeof input === "string") return input;
  if (input instanceof URL) return input.href;
  return input.url;
}

it("creates a meeting type and refreshes the list", async () => {
  let meetingTypes: MeetingType[] = [];
  const fetchMock = vi.fn<typeof fetch>().mockImplementation((input, init) => {
    const method = init?.method ?? "GET";
    const url = requestUrl(input);
    if (url.endsWith("/api/v1/meeting-types") && method === "POST") {
      const created: MeetingType = {
        id: "intro",
        name: "Вводный звонок",
        description: "Знакомство с командой",
        durationMinutes: MeetingDuration.NUMBER_60,
      };
      meetingTypes = [created];
      return Promise.resolve(jsonResponse(created, 201));
    }
    if (url.endsWith("/api/v1/meeting-types") && method === "GET") {
      return Promise.resolve(jsonResponse({ items: meetingTypes }));
    }
    return Promise.resolve(jsonResponse({ message: "not found" }, 404));
  });
  vi.stubGlobal("fetch", fetchMock);

  renderApp("/owner");

  expect(await screen.findByText("Пока нет типов встреч")).toBeInTheDocument();

  fireEvent.change(screen.getByLabelText(/ID/), { target: { value: "intro" } });
  fireEvent.change(screen.getByLabelText(/Название/), {
    target: { value: "Вводный звонок" },
  });
  fireEvent.change(screen.getByLabelText(/Описание/), {
    target: { value: "Знакомство с командой" },
  });
  fireEvent.change(screen.getByLabelText(/Длительность/), {
    target: { value: String(MeetingDuration.NUMBER_60) },
  });
  fireEvent.click(screen.getByRole("button", { name: "Создать тип встречи" }));

  expect(
    await screen.findByRole("heading", { name: "Вводный звонок" }),
  ).toBeInTheDocument();
  expect(screen.getByLabelText<HTMLInputElement>(/Название/).value).toBe("");

  const postCall = fetchMock.mock.calls.find(
    ([, init]) => init?.method === "POST",
  );
  expect(postCall).toBeDefined();
  const requestBody = postCall?.[1]?.body;
  if (typeof requestBody !== "string") {
    throw new Error("request body is not a string");
  }
  const body: unknown = JSON.parse(requestBody);
  expect(body).toMatchObject({
    id: "intro",
    name: "Вводный звонок",
    description: "Знакомство с командой",
    durationMinutes: MeetingDuration.NUMBER_60,
  });
});

it("shows the server validation message and keeps the list unchanged", async () => {
  const fetchMock = vi.fn<typeof fetch>().mockImplementation((input, init) => {
    const method = init?.method ?? "GET";
    if (
      requestUrl(input).endsWith("/api/v1/meeting-types") &&
      method === "POST"
    ) {
      return Promise.resolve(
        jsonResponse({ message: "Название не может быть пустым" }, 400),
      );
    }
    return Promise.resolve(jsonResponse({ items: [] }));
  });
  vi.stubGlobal("fetch", fetchMock);

  renderApp("/owner");

  expect(await screen.findByText("Пока нет типов встреч")).toBeInTheDocument();

  fireEvent.change(screen.getByLabelText(/ID/), { target: { value: "intro" } });
  fireEvent.change(screen.getByLabelText(/Название/), {
    target: { value: "   " },
  });
  fireEvent.change(screen.getByLabelText(/Описание/), {
    target: { value: "Знакомство с командой" },
  });
  fireEvent.click(screen.getByRole("button", { name: "Создать тип встречи" }));

  expect(
    await screen.findByText("Название не может быть пустым"),
  ).toBeInTheDocument();
  expect(screen.getByLabelText<HTMLInputElement>(/Название/).value).toBe("   ");
  expect(screen.getByText("Пока нет типов встреч")).toBeInTheDocument();
});
