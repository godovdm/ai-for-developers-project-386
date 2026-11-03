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
