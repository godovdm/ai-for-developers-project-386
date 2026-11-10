import {
  Alert,
  Button,
  Container,
  Loader,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useMeetingTypes } from "../hooks/useMeetingTypes";
import { AppLink } from "../router";

interface MeetingTypePageProps {
  id: string;
}

export function MeetingTypePage({ id }: MeetingTypePageProps) {
  const { state, reload } = useMeetingTypes();
  const meetingType =
    state.status === "ready"
      ? state.meetingTypes.find((item) => item.id === id)
      : undefined;

  return (
    <Container size="md" py="xl">
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
          <Stack gap="sm">
            <Title order={1}>{meetingType.name}</Title>
            <Text c="dimmed">{meetingType.description}</Text>
            <Text>Длительность: {meetingType.durationMinutes} минут</Text>
            <Alert title="Следующий шаг">
              Выбор даты и времени появится на следующем шаге.
            </Alert>
          </Stack>
        )}
      </Stack>
    </Container>
  );
}
