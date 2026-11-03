import { Alert, Container, Loader, Title } from "@mantine/core";
import { useEffect, useState } from "react";
import { DefaultApi } from "./api/generated/apis/DefaultApi";
import { Configuration } from "./api/generated/runtime";

const api = new DefaultApi(new Configuration({ basePath: "" }));

export function App() {
  const [status, setStatus] = useState<"loading" | "ok" | "error">("loading");

  useEffect(() => {
    let active = true;
    void api.getHealth().then(
      () => {
        if (active) setStatus("ok");
      },
      () => {
        if (active) setStatus("error");
      },
    );
    return () => {
      active = false;
    };
  }, []);

  return (
    <Container size="sm" py="xl">
      <Title>Календарь встреч</Title>
      {status === "loading" && <Loader aria-label="Проверка сервера" />}
      {status === "ok" && <Alert color="green">Сервер работает</Alert>}
      {status === "error" && <Alert color="red">Сервер недоступен</Alert>}
    </Container>
  );
}
