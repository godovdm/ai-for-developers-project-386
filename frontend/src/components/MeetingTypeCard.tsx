import { Button, Group, Paper, Stack, Text, Title } from "@mantine/core";
import type { MeetingType } from "../api/generated/models/MeetingType";
import { navigate } from "../router";

interface MeetingTypeCardProps {
  meetingType: MeetingType;
}

export function MeetingTypeCard({ meetingType }: MeetingTypeCardProps) {
  return (
    <Paper withBorder p="lg" radius="md">
      <Stack gap="sm" h="100%">
        <Group justify="space-between" align="flex-start" wrap="nowrap">
          <Title order={3} size="h4">
            {meetingType.name}
          </Title>
          <Text size="sm" c="dimmed" style={{ whiteSpace: "nowrap" }}>
            {meetingType.durationMinutes} минут
          </Text>
        </Group>
        <Text size="sm" c="dimmed" style={{ flexGrow: 1 }}>
          {meetingType.description}
        </Text>
        <Button
          variant="default"
          onClick={() => {
            navigate(`/meeting-types/${encodeURIComponent(meetingType.id)}`);
          }}
        >
          Выбрать время
        </Button>
      </Stack>
    </Paper>
  );
}
