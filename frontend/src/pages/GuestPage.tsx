import {
  Alert,
  Box,
  Button,
  Center,
  Container,
  Loader,
  Paper,
  Stack,
  Text,
  Title,
  VisuallyHidden,
} from "@mantine/core";
import { MeetingTypeCard } from "../components/MeetingTypeCard";
import { useMeetingTypes } from "../hooks/useMeetingTypes";

export function GuestPage() {
  const { state, reload } = useMeetingTypes();
  return (
    <Container size="lg" py={{ base: "lg", sm: 48 }}>
      <Stack gap={48}>
        <VisuallyHidden>
          <Title order={1}>Выберите время для встречи</Title>
        </VisuallyHidden>
        <Paper className="public-profile" withBorder radius="xl">
          <Box className="profile-avatar" aria-hidden="true">
            Д
          </Box>
          <Title order={2}>Дмитрий</Title>
          <Text c="dimmed">Выберите подходящий формат встречи</Text>
        </Paper>

        <Box
          id="meeting-types"
          component="section"
          aria-labelledby="meeting-types-title"
        >
          <Title id="meeting-types-title" order={2} mb="md" size="h3">
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
            <Paper className="meeting-type-list" withBorder radius="lg">
              {state.meetingTypes.map((meetingType) => (
                <MeetingTypeCard
                  key={meetingType.id}
                  meetingType={meetingType}
                />
              ))}
            </Paper>
          )}
        </Box>
      </Stack>
    </Container>
  );
}
