import { useCallback, useEffect, useState } from "react";
import { api } from "../api/client";
import type { MeetingType } from "../api/generated/models/MeetingType";

type MeetingTypesState =
  | { status: "loading" }
  | { status: "ready"; meetingTypes: MeetingType[] }
  | { status: "error" };

export function useMeetingTypes() {
  const [state, setState] = useState<MeetingTypesState>({ status: "loading" });
  const [request, setRequest] = useState(0);

  useEffect(() => {
    let active = true;
    void api.listMeetingTypes().then(
      (meetingTypes) => {
        if (active) setState({ status: "ready", meetingTypes });
      },
      () => {
        if (active) setState({ status: "error" });
      },
    );
    return () => {
      active = false;
    };
  }, [request]);

  const reload = useCallback(() => {
    setState({ status: "loading" });
    setRequest((value) => value + 1);
  }, []);

  return { state, reload };
}
