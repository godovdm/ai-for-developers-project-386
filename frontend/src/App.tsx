import { Box, Button, Container, Stack, Text, Title } from "@mantine/core";
import { AppHeader } from "./components/AppHeader";
import { usePath } from "./router";
import { GuestPage } from "./pages/GuestPage";
import { MeetingTypePage } from "./pages/MeetingTypePage";
import { OwnerPage } from "./pages/OwnerPage";

function renderPage(path: string) {
  if (path === "/" || path === "") return <GuestPage />;
  if (path === "/owner") return <OwnerPage />;
  const match = /^\/meeting-types\/([^/]+)$/.exec(path);
  const id = match?.[1];
  if (id !== undefined) {
    const decodedID = safeDecode(id);
    if (decodedID !== undefined) return <MeetingTypePage id={decodedID} />;
  }
  return <NotFoundPage />;
}

function safeDecode(value: string): string | undefined {
  try {
    return decodeURIComponent(value);
  } catch {
    return undefined;
  }
}

function NotFoundPage() {
  return (
    <Container size="sm" py="xl">
      <Stack gap="md">
        <Title order={1}>Страница не найдена</Title>
        <Text c="dimmed">
          Проверьте адрес или вернитесь на главную страницу.
        </Text>
        <Button component="a" href="/" variant="default" w="fit-content">
          На главную
        </Button>
      </Stack>
    </Container>
  );
}

export function App() {
  const path = usePath();

  return (
    <Box style={{ minHeight: "100vh" }}>
      <AppHeader path={path} />
      <main>{renderPage(path)}</main>
    </Box>
  );
}
