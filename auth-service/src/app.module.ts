import { Module } from "@nestjs/common";
import { TypeOrmModule } from "@nestjs/typeorm";
import { JwtModule } from "@nestjs/jwt";
import { AuthController } from "./controller/auth.controller";
import { AuthService } from "./service/auth.service";
import { User } from "./entities/user.entity";
import { RefreshToken } from "./entities/refresh-token.entity";

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: "postgres",
      url: process.env.DATABASE_URL,
      entities: [User, RefreshToken],
      synchronize: false, // Use migrations in production
      logging: process.env.NODE_ENV !== "production",
    }),
    TypeOrmModule.forFeature([User, RefreshToken]),
    JwtModule.register({
      secret: process.env.JWT_SECRET || "your-secret-key",
      signOptions: { expiresIn: process.env.JWT_EXPIRES_IN || "15m" },
    }),
  ],
  controllers: [AuthController],
  providers: [AuthService],
})
export class AppModule {}
