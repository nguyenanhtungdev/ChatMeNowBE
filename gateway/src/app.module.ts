import { Module } from "@nestjs/common";
import { ThrottlerModule } from "@nestjs/throttler";
import { JwtModule } from "@nestjs/jwt";
import { HttpModule } from "@nestjs/axios";
import { AppController } from "./controllers/app.controller";
import { AuthProxyController } from "./controllers/auth-proxy.controller";
import { BlogProxyController } from "./controllers/blog-proxy.controller";
import { ChatProxyController } from "./controllers/chat-proxy.controller";
import { JwtAuthGuard } from "./guards/jwt-auth.guard";
import { ProxyService } from "./services/proxy.service";

@Module({
  imports: [
    // Rate limiting
    ThrottlerModule.forRoot([
      {
        ttl: parseInt(process.env.RATE_LIMIT_TTL || "60", 10) * 1000,
        limit: parseInt(process.env.RATE_LIMIT_MAX || "100", 10),
      },
    ]),
    // JWT
    JwtModule.register({
      secret: process.env.JWT_SECRET || "your-secret-key",
      signOptions: { expiresIn: "15m" },
    }),
    // HTTP client for proxying
    HttpModule.register({
      timeout: 30000,
      maxRedirects: 5,
    }),
  ],
  controllers: [
    AppController,
    AuthProxyController,
    BlogProxyController,
    ChatProxyController,
  ],
  providers: [JwtAuthGuard, ProxyService],
})
export class AppModule {}
