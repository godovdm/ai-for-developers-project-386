import {
  ActionIcon,
  Box,
  Container,
  Group,
  useComputedColorScheme,
  useMantineColorScheme,
} from "@mantine/core";
import { AppLink } from "../router";

interface AppHeaderProps {
  path: string;
}

export function AppHeader({ path }: AppHeaderProps) {
  const { toggleColorScheme } = useMantineColorScheme();
  const colorScheme = useComputedColorScheme("light", {
    getInitialValueInEffect: false,
  });
  const ownerView = path === "/owner";

  return (
    <Box
      component="header"
      py="md"
      style={{ borderBottom: "1px solid var(--mantine-color-default-border)" }}
    >
      <Container size="lg">
        <Group justify="space-between" align="center">
          <AppLink
            to="/"
            fw={600}
            style={{ color: "var(--mantine-color-text)" }}
            underline="never"
          >
            Calendar
          </AppLink>
          <Group gap="lg">
            <AppLink
              to={ownerView ? "/" : "/owner"}
              style={{ color: "var(--mantine-color-text)" }}
              underline="hover"
            >
              {ownerView ? "Гостю" : "Владельцу"}
            </AppLink>
            <ActionIcon
              variant="default"
              aria-label="Переключить тему"
              onClick={() => {
                toggleColorScheme();
              }}
            >
              <span aria-hidden="true">
                {colorScheme === "dark" ? "☀" : "☾"}
              </span>
            </ActionIcon>
          </Group>
        </Group>
      </Container>
    </Box>
  );
}
