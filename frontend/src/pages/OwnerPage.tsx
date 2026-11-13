import {
  Alert,
  Box,
  Button,
  Center,
  Container,
  Group,
  Loader,
  NativeSelect,
  Paper,
  SimpleGrid,
  Stack,
  Text,
  Textarea,
  TextInput,
  Title,
} from "@mantine/core";
import { useEffect, useState, type SyntheticEvent } from "react";
import { api, readErrorMessage } from "../api/client";
import type { UpcomingBooking } from "../api/generated/models/UpcomingBooking";
import { MeetingDuration } from "../api/generated/models/MeetingDuration";
import { useMeetingTypes } from "../hooks/useMeetingTypes";

const DURATION_OPTIONS: { value: string; label: string }[] = [
  { value: String(MeetingDuration.NUMBER_30), label: "30 минут" },
  { value: String(MeetingDuration.NUMBER_60), label: "60 минут" },
  { value: String(MeetingDuration.NUMBER_90), label: "90 минут" },
  { value: String(MeetingDuration.NUMBER_120), label: "120 минут" },
];

function parseDuration(value: string): MeetingDuration {
  switch (Number(value)) {
    case MeetingDuration.NUMBER_30:
      return MeetingDuration.NUMBER_30;
    case MeetingDuration.NUMBER_90:
      return MeetingDuration.NUMBER_90;
    case MeetingDuration.NUMBER_120:
      return MeetingDuration.NUMBER_120;
    default:
      return MeetingDuration.NUMBER_60;
  }
}

export function OwnerPage() {
  const { state, reload } = useMeetingTypes();
  const [bookings, setBookings] = useState<UpcomingBooking[] | null>(null);
  const [bookingError, setBookingError] = useState(false);
  useEffect(() => {
    let active = true;
    void api.listUpcomingBookings().then(
      (result) => {
        if (active) setBookings(result.items);
      },
      () => {
        if (active) setBookingError(true);
      },
    );
    return () => {
      active = false;
    };
  }, []);
  const [id, setId] = useState("");
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [duration, setDuration] = useState<MeetingDuration>(
    MeetingDuration.NUMBER_60,
  );
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (event: SyntheticEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      await api.createMeetingTypeV1({
        createMeetingTypeV1Request: {
          id,
          name,
          description,
          durationMinutes: duration,
        },
      });
      setId("");
      setName("");
      setDescription("");
      setDuration(MeetingDuration.NUMBER_60);
      reload();
    } catch (cause) {
      setError(await readErrorMessage(cause, "Не удалось создать тип встречи"));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Container size="lg" py="xl">
      <Stack gap="xl">
        <Paper withBorder p="lg" radius="md">
          <form
            onSubmit={(event) => {
              void handleSubmit(event);
            }}
          >
            <Stack gap="md">
              <Title order={1} size="h2">
                Новый тип встречи
              </Title>
              <TextInput
                label="ID"
                value={id}
                maxLength={64}
                required
                pattern="[a-z0-9]([a-z0-9-]{0,62}[a-z0-9])?"
                onChange={(event) => {
                  setId(event.currentTarget.value);
                }}
              />
              <TextInput
                label="Название"
                value={name}
                maxLength={100}
                required
                onChange={(event) => {
                  setName(event.currentTarget.value);
                }}
              />
              <Textarea
                label="Описание"
                value={description}
                maxLength={2000}
                minRows={3}
                onChange={(event) => {
                  setDescription(event.currentTarget.value);
                }}
              />
              <NativeSelect
                label="Длительность"
                value={String(duration)}
                data={DURATION_OPTIONS}
                onChange={(event) => {
                  setDuration(parseDuration(event.currentTarget.value));
                }}
              />
              {error !== null && <Alert color="red">{error}</Alert>}
              <Button
                className="app-primary"
                type="submit"
                loading={submitting}
                w="fit-content"
              >
                Создать тип встречи
              </Button>
            </Stack>
          </form>
        </Paper>

        <Box component="section" aria-label="Предстоящие записи">
          <Title order={2} mb="md">
            Предстоящие записи
          </Title>
          {bookingError && (
            <Alert color="red">Не удалось загрузить записи</Alert>
          )}
          {!bookingError && bookings === null && (
            <Loader aria-label="Загрузка записей" />
          )}
          {bookings?.length === 0 && (
            <Text c="dimmed">Пока нет предстоящих записей</Text>
          )}
          {bookings?.map((booking) => (
            <Paper
              key={booking.start + booking.guestEmail}
              withBorder
              p="md"
              mb="sm"
            >
              <Text fw={600}>{booking.meetingType.name}</Text>
              <Text>
                {booking.start.slice(0, 10)} {booking.start.slice(11, 16)}{" "}
                (Москва)
              </Text>
              <Text>{booking.guestName}</Text>
              <Text>{booking.guestEmail}</Text>
            </Paper>
          ))}
        </Box>
        <Box>
          <Title order={2} mb="md">
            Созданные типы встреч
          </Title>
          {state.status === "loading" && (
            <Center py="xl">
              <Loader aria-label="Загрузка типов встреч" />
            </Center>
          )}
          {state.status === "error" && (
            <Alert color="red" title="Не удалось загрузить типы встреч">
              <Button variant="default" onClick={reload}>
                Попробовать снова
              </Button>
            </Alert>
          )}
          {state.status === "ready" && state.meetingTypes.length === 0 && (
            <Text c="dimmed">Пока нет типов встреч</Text>
          )}
          {state.status === "ready" && state.meetingTypes.length > 0 && (
            <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="md">
              {state.meetingTypes.map((meetingType) => (
                <Paper key={meetingType.id} withBorder p="lg" radius="md">
                  <Stack gap="xs">
                    <Group
                      justify="space-between"
                      align="flex-start"
                      wrap="nowrap"
                    >
                      <Title order={3} size="h4">
                        {meetingType.name}
                      </Title>
                      <Text
                        size="sm"
                        c="dimmed"
                        style={{ whiteSpace: "nowrap" }}
                      >
                        {meetingType.durationMinutes} минут
                      </Text>
                    </Group>
                    <Text size="sm" c="dimmed">
                      {meetingType.description}
                    </Text>
                  </Stack>
                </Paper>
              ))}
            </SimpleGrid>
          )}
        </Box>
      </Stack>
    </Container>
  );
}
