import { Box, Stack, Text, Title } from "@mantine/core";
import type { MeetingType } from "../api/generated/models/MeetingType";
import { navigate } from "../router";

interface MeetingTypeCardProps {
  meetingType: MeetingType;
}

export function MeetingTypeCard({ meetingType }: MeetingTypeCardProps) {
  const titleID = `meeting-type-${meetingType.id}`;

  return (
    <Box className="meeting-type-card">
      <Stack className="meeting-type-content" gap="xs" align="flex-start">
        <Title id={titleID} order={3} size="h4">
          {meetingType.name}
        </Title>
        {meetingType.description && (
          <Text size="sm" c="dimmed">
            {meetingType.description}
          </Text>
        )}
        <Text className="duration-badge" size="sm">
          <span aria-hidden="true">◷</span> {meetingType.durationMinutes} минут
        </Text>
      </Stack>
      <span className="meeting-type-arrow" aria-hidden="true">
        →
      </span>
      <button
        type="button"
        className="meeting-card-action"
        aria-label="Выбрать время"
        aria-describedby={titleID}
        onClick={() => {
          navigate(`/meeting-types/${encodeURIComponent(meetingType.id)}`);
        }}
      />
    </Box>
  );
}
