import { NestFactory } from "@nestjs/core";
import { ValidationPipe, Logger } from "@nestjs/common";
import { AppModule } from "./app.module";

async function bootstrap() {
  const logger = new Logger("Gateway");
  const app = await NestFactory.create(AppModule);

  // Enable CORS
  app.enableCors({
    origin: true,
    credentials: true,
  });

  // Set global API prefix
  app.setGlobalPrefix("api");

  // Global validation pipe
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      transform: true,
    }),
  );

  // Add request ID and logging
  app.use((req, res, next) => {
    req.id = Math.random().toString(36).substring(7);
    logger.log(`[${req.id}] ${req.method} ${req.url}`);
    next();
  });

  const port = process.env.PORT || 3000;
  await app.listen(port);
  logger.log(`Gateway running on http://localhost:${port}`);
}

bootstrap();
