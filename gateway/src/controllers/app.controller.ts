import { Controller, Get } from "@nestjs/common";

@Controller()
export class AppController {
  @Get("health")
  healthCheck() {
    return {
      status: "ok",
      service: "gateway",
      timestamp: new Date().toISOString(),
    };
  }

  @Get()
  root() {
    return {
      name: "ChatMeNow API Gateway",
      version: "1.0.0",
      endpoints: {
        health: "GET /health",
        auth: "POST /api/auth/*",
        blog: "/api/blog/*",
        chat: "/api/chat/*",
      },
    };
  }
}
