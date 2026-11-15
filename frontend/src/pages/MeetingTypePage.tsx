import {
  Alert,
  Box,
  Button,
  Container,
  Group,
  Loader,
  Paper,
  SimpleGrid,
  Stack,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { useEffect, useState } from "react";
import { api } from "../api/client";
import { ResponseError, type Availability } from "../api/generated";
import { useMeetingTypes } from "../hooks/useMeetingTypes";
import { AppLink } from "../router";

interface MeetingTypePageProps {
  id: string;
}

export function MeetingTypePage({ id }: MeetingTypePageProps) {
  const [day, setDay] = useState(() =>
    new Date().toLocaleDateString("en-CA", { timeZone: "Europe/Moscow" }),
  );
  useEffect(() => {
    let timer: number;
    const updateDay = () => {
      const now = new Date();
      setDay(now.toLocaleDateString("en-CA", { timeZone: "Europe/Moscow" }));
      window.clearTimeout(timer);
      const moscowTime = now.getTime() + 3 * 60 * 60_000;
      timer = window.setTimeout(
        updateDay,
        Math.min(30_000, 86_400_000 - (moscowTime % 86_400_000)),
      );
    };
    updateDay();
    window.addEventListener("focus", updateDay);
    window.addEventListener("pageshow", updateDay);
    document.addEventListener("visibilitychange", updateDay);
    return () => {
      window.clearTimeout(timer);
      window.removeEventListener("focus", updateDay);
      window.removeEventListener("pageshow", updateDay);
      document.removeEventListener("visibilitychange", updateDay);
    };
  }, []);
  return <MeetingTypeContent key={`${id}:${day}`} id={id} />;
}

function MeetingTypeContent({ id }: MeetingTypePageProps) {
  const { state, reload } = useMeetingTypes();
  const [availability, setAvailability] = useState<Availability>();
  const [selectedDate, setSelectedDate] = useState("");
  const [selectedSlot, setSelectedSlot] = useState("");
  const [availabilityError, setAvailabilityError] = useState(false);
  const [guestName, setGuestName] = useState("");
  const [guestEmail, setGuestEmail] = useState("");
  const [bookingError, setBookingError] = useState("");
  const [confirmed, setConfirmed] = useState("");
  const [submitting, setSubmitting] = useState(false);
  async function refreshAvailability() {
    try {
      setAvailability(await api.getAvailability({ id }));
      setAvailabilityError(false);
    } catch {
      setAvailabilityError(true);
    }
  }
  async function book() {
    setSubmitting(true);
    setBookingError("");
    try {
      const result = await api.createBooking({
        createBookingRequest: {
          meetingTypeId: id,
          start: selectedSlot,
          guestName,
          guestEmail,
        },
      });
      setConfirmed(
        `${result.meetingType.name}: ${result.start.slice(0, 10)} ${result.start.slice(11, 16)} (Москва)`,
      );
      setSelectedSlot("");
      await refreshAvailability();
    } catch (error) {
      if (error instanceof ResponseError && error.response.status === 409) {
        setBookingError("Это время уже занято. Выберите другое время.");
        setSelectedSlot("");
        await refreshAvailability();
      } else {
        setBookingError(
          "Не удалось записаться. Проверьте имя и email и попробуйте снова.",
        );
      }
    } finally {
      setSubmitting(false);
    }
  }
  useEffect(() => {
    let active = true;
    api
      .getAvailability({ id })
      .then((result) => {
        if (active) {
          setAvailability(result);
          setSelectedDate(result.days[0]?.date ?? "");
        }
      })
      .catch(() => {
        if (active) setAvailabilityError(true);
      });
    return () => {
      active = false;
    };
  }, [id]);
  const meetingType =
    state.status === "ready"
      ? state.meetingTypes.find((item) => item.id === id)
      : undefined;

  return (
    <Container size="xl" py={{ base: "md", sm: "xl" }}>
      <Stack gap="lg">
        <AppLink to="/">← Все типы встреч</AppLink>
        {state.status === "loading" && (
          <Loader aria-label="Загрузка типа встречи" />
        )}
        {state.status === "error" && (
          <Alert color="red" title="Не удалось загрузить типы встреч">
            <Button variant="default" onClick={reload}>
              Попробовать снова
            </Button>
          </Alert>
        )}
        {state.status === "ready" && meetingType === undefined && (
          <Title order={1}>Тип встречи не найден</Title>
        )}
        {meetingType !== undefined && (
          <Paper className="booking-shell" withBorder radius="lg">
            <Box className="booking-details">
              <Stack gap="lg">
                <Box className="booking-avatar" aria-hidden="true">
                  Д
                </Box>
                <Stack gap="xs">
                  <Text c="dimmed" fw={600}>
                    Дмитрий
                  </Text>
                  <Title order={1} size="h2">
                    {meetingType.name}
                  </Title>
                </Stack>
                {meetingType.description && (
                  <Text c="dimmed">{meetingType.description}</Text>
                )}
                <Text className="booking-meta">
                  <span aria-hidden="true">◷</span>{" "}
                  {meetingType.durationMinutes} минут
                </Text>
                <Text className="booking-meta">
                  <span aria-hidden="true">◎</span> Europe/Moscow
                </Text>
              </Stack>
            </Box>

            <Box className="booking-calendar">
              <Group justify="space-between" mb="xl">
                <Title order={2} size="h3">
                  Доступные даты
                </Title>
                <Text c="dimmed" size="sm">
                  Время Москвы
                </Text>
              </Group>
              {availabilityError && (
                <Alert color="red" mb="md">
                  Не удалось загрузить доступное время
                </Alert>
              )}
              {!availability && !availabilityError && (
                <Loader aria-label="Загрузка доступного времени" />
              )}
              {availability && (
                <SimpleGrid cols={7} spacing="xs" verticalSpacing="xs">
                  {availability.days.map((availableDay) => {
                    const date = new Date(`${availableDay.date}T12:00:00`);
                    return (
                      <Button
                        key={availableDay.date}
                        className="calendar-day"
                        aria-label={availableDay.date}
                        variant={
                          selectedDate === availableDay.date
                            ? "filled"
                            : "default"
                        }
                        disabled={availableDay.slots.length === 0}
                        onClick={() => {
                          setSelectedDate(availableDay.date);
                          setSelectedSlot("");
                        }}
                      >
                        <span>
                          {date.toLocaleDateString("ru-RU", {
                            weekday: "short",
                          })}
                        </span>
                        <strong>{date.getDate()}</strong>
                      </Button>
                    );
                  })}
                </SimpleGrid>
              )}
            </Box>

            <Box className="booking-times">
              <Title order={2} size="h3" mb="lg">
                {selectedDate
                  ? new Date(`${selectedDate}T12:00:00`).toLocaleDateString(
                      "ru-RU",
                      {
                        weekday: "long",
                        day: "numeric",
                        month: "long",
                      },
                    )
                  : "Выберите дату"}
              </Title>
              <Stack gap="sm">
                {availability?.days
                  .find((availableDay) => availableDay.date === selectedDate)
                  ?.slots.map((slot) => (
                    <Button
                      key={slot}
                      className="time-slot"
                      variant={selectedSlot === slot ? "filled" : "default"}
                      onClick={() => {
                        setConfirmed("");
                        setSelectedSlot(slot);
                      }}
                    >
                      <span aria-hidden="true" className="availability-dot" />
                      {slot.slice(11, 16)}
                    </Button>
                  ))}
                {availability?.days.find(
                  (availableDay) => availableDay.date === selectedDate,
                )?.slots.length === 0 && (
                  <Text c="dimmed">Свободных слотов нет</Text>
                )}
                {bookingError && <Alert color="red">{bookingError}</Alert>}
                {confirmed && (
                  <Alert color="green" title="Запись подтверждена">
                    {confirmed}
                  </Alert>
                )}
                {selectedSlot && !confirmed && (
                  <form
                    className="booking-form"
                    onSubmit={(event) => {
                      event.preventDefault();
                      void book();
                    }}
                  >
                    <Stack>
                      <Text size="sm" fw={600}>
                        Выбрано: {selectedSlot.slice(0, 10)}{" "}
                        {selectedSlot.slice(11, 16)} (Москва)
                      </Text>
                      <TextInput
                        label="Имя"
                        value={guestName}
                        onChange={(event) => {
                          setGuestName(event.currentTarget.value);
                        }}
                        required
                      />
                      <TextInput
                        label="Email"
                        type="email"
                        value={guestEmail}
                        onChange={(event) => {
                          setGuestEmail(event.currentTarget.value);
                        }}
                        required
                      />
                      <Button
                        className="app-primary"
                        type="submit"
                        loading={submitting}
                      >
                        Записаться
                      </Button>
                    </Stack>
                  </form>
                )}
              </Stack>
            </Box>
          </Paper>
        )}
      </Stack>
    </Container>
  );
}
