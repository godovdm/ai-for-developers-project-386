import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

using TypeSpec.Http;

@service(#{ title: "Meeting calendar API" })
@useAuth(NoAuth)
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

@doc("Meeting type creation using an owner-supplied identifier.")
model CreateMeetingTypeV1Request {
  @minLength(1)
  @maxLength(64)
  @pattern("^[a-z0-9]([a-z0-9-]{0,62}[a-z0-9])?$")
  id: string;

  @minLength(1)
  @maxLength(100)
  name: string;

  @maxLength(2000)
  description: string;

  durationMinutes: MeetingDuration;
}

model MeetingTypeV1 {
  id: string;
  name: string;
  description: string;
  durationMinutes: MeetingDuration;
}

model MeetingTypeListV1 {
  items: MeetingTypeV1[];
}

@error
model ApiErrorV1 {
  code: string;
  message: string;
}

@summary("List meeting types using API v1")
@route("/api/v1/meeting-types")
@get
op listMeetingTypesV1():
  | MeetingTypeListV1
  | {
      @statusCode statusCode: 400;
      @body body: ApiErrorV1;
    };

@summary("Create a meeting type using API v1")
@route("/api/v1/meeting-types")
@post
op createMeetingTypeV1(@body body: CreateMeetingTypeV1Request):
  | {
      @statusCode statusCode: 201;
      @body body: MeetingTypeV1;
    }
  | {
      @statusCode statusCode: 400;
      @body body: ApiErrorV1;
    }
  | {
      @statusCode statusCode: 409;
      @body body: ApiErrorV1;
    };

@summary("Get a meeting type using API v1")
@route("/api/v1/meeting-types/{id}")
@get
op getMeetingTypeV1(@path id: string):
  | MeetingTypeV1
  | {
      @statusCode statusCode: 404;
      @body body: ApiErrorV1;
    };

@doc("Available starts for one calendar day in the owner's time zone.")
model AvailabilityDay {
  date: string;
  slots: string[];
}

@doc("The current 14-day booking window.")
model Availability {
  days: AvailabilityDay[];
}

@summary("List available starts for a meeting type")
@route("/api/v1/meeting-types/{id}/availability")
@get
op getAvailability(@path id: string, @query from?: string, @query to?: string):
  | Availability
  | {
      @statusCode statusCode: 404;
      @body body: ApiErrorV1;
    };
