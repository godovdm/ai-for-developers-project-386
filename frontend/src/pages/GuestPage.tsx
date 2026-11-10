import {
  Alert,
  Box,
  Button,
  Center,
  Container,
  List,
  Loader,
  Paper,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useRef } from "react";
import { MeetingTypeCard } from "../components/MeetingTypeCard";
import { useMeetingTypes } from "../hooks/useMeetingTypes";

export function GuestPage() {
  const { state, reload } = useMeetingTypes();
  const listRef = useRef<HTMLDivElement>(null);

  return (
    <Container size="lg" py="xl">
      <SimpleGrid cols={{ base: 1, md: 2 }} spacing="xl" verticalSpacing="xl">
        <Stack gap="md" justify="center">
          <Title order={1}>Выберите время для встречи</Title>
          <Text size="lg" c="dimmed">
            Выберите тип встречи и удобное время. Подтверждение появится сразу
            после записи
          </Text>
          <Button
            className="app-primary"
            size="md"
            w="fit-content"
            onClick={() => {
              listRef.current?.scrollIntoView({
                behavior: "smooth",
                block: "start",
              });
            }}
          >
            Записаться
          </Button>
        </Stack>
        <Paper withBorder p="lg" radius="md">
          <Title order={2} size="h3" mb="md">
            Возможности
          </Title>
          <List spacing="sm" size="sm">
            <List.Item>Выбор типа встречи</List.Item>
            <List.Item>Просмотр свободного времени</List.Item>
            <List.Item>Запись без регистрации</List.Item>
          </List>
        </Paper>
      </SimpleGrid>

      <Box
        ref={listRef}
        id="meeting-types"
        mt={48}
        style={{ scrollMarginTop: 24 }}
      >
        <Title order={2} mb="md">
          Типы встреч
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
          <Text c="dimmed">Пока нет доступных встреч</Text>
        )}
        {state.status === "ready" && state.meetingTypes.length > 0 && (
          <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="md">
            {state.meetingTypes.map((meetingType) => (
              <MeetingTypeCard key={meetingType.id} meetingType={meetingType} />
            ))}
          </SimpleGrid>
        )}
      </Box>
    </Container>
  );
}
