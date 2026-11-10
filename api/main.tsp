import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

using TypeSpec.Http;

@service(#{ title: "Meeting calendar API" })
@OpenAPI.info(#{ version: "0.1.0" })
@server("/", "Same origin as the application")
namespace Calendar;

@doc("Health status of the server.")
model HealthStatus {
  @doc("The server is responding.")
  status: "ok";
}

@summary("Check that the server is running")
@route("/api/health")
@get
op getHealth(): HealthStatus;

@doc("Allowed meeting duration in minutes.")
union MeetingDuration {
  @doc("30 minutes")
  thirty: 30,

  @doc("60 minutes")
  sixty: 60,

  @doc("90 minutes")
  ninety: 90,

  @doc("120 minutes")
  oneHundredTwenty: 120,
}

@doc("A meeting type offered by the calendar owner.")
model MeetingType {
  @doc("Unique meeting type identifier.")
  id: string;

  @doc("Meeting type name.")
  @maxLength(100)
  name: string;

  @doc("Meeting type description.")
  @maxLength(1000)
  description: string;

  @doc("Meeting duration in minutes.")
  durationMinutes: MeetingDuration;
}

@doc("Data required to create a meeting type.")
model CreateMeetingTypeRequest {
  @doc("Meeting type name.")
  @maxLength(100)
  name: string;

  @doc("Meeting type description.")
  @maxLength(1000)
  description: string;

  @doc("Meeting duration in minutes.")
  durationMinutes: MeetingDuration;
}

@error
@doc("API error with a human-readable message.")
model ApiError {
  @doc("HTTP status code.")
  @statusCode
  statusCode: 400;

  @doc("Human-readable error message.")
  message: string;
}

@summary("List meeting types")
@route("/api/meeting-types")
@get
op listMeetingTypes(): MeetingType[];

@summary("Create a meeting type")
@route("/api/meeting-types")
@post
op createMeetingType(@body body: CreateMeetingTypeRequest):
  | {
      @statusCode statusCode: 201;
      @body body: MeetingType;
    }
  | ApiError;
